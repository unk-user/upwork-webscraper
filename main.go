package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const baseUrl = "https://www.upwork.com/nx/search/jobs/"

type Response struct {
	StatusCode int `json:"statusCode"`
	Length     int `json:"length"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	keywords := strings.Fields(event.QueryStringParameters["keywords"]) // Use strings.Fields to split the keywords into an array
	categoryId := event.QueryStringParameters["categoryId"]
	apiEndpoint := event.QueryStringParameters["apiEndpoint"]

	html, err := GetHTML(baseUrl + "?category2_uid=" + categoryId + "&per_page=20" + "&q=%28" + strings.Join(keywords, "%20OR%20") + "%29")
	if err != nil {
		log.Println("Error getting HTML", err.Error())
		return Response{}, err
	}

	jobs, err := ProcessHTML(html)
	if err != nil {
		log.Println("Error parsing HTML", err.Error())
		return Response{}, err
	}

	err = sendJobs(jobs, apiEndpoint)
	if err != nil {
		log.Println("Error sending jobs", err.Error())
		return Response{}, err
	}

	log.Println("Scraping job finished successfully")
	return Response{StatusCode: 200, Length: len(jobs)}, nil
}

func main() {
	lambda.Start(handler)
}

func sendJobs(jobs []Job, apiEndpoint string) error {
	if apiEndpoint == "" {
		log.Println("API endpoint not specified")
		return nil
	}

	data, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
