package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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

	html, err := GetHTML(baseUrl + "?category2_uid=" + categoryId + "&per_page=20" + "&q=%28" + strings.Join(keywords, "%20OR%20") + "%29")
	if err != nil {
		log.Println(err)
		return Response{}, err
	}

	jobs, err := ProcessHTML(html)
	if err != nil {
		log.Println(err)
		return Response{}, err
	}

	err = sendJobs(jobs)
	if err != nil {
		log.Println(err)
		return Response{}, err
	}

	log.Println("Done")
	return Response{StatusCode: 200, Length: len(jobs)}, nil
}

func main() {
	lambda.Start(handler)
}

func sendJobs(jobs []Job) error {
	const apiEndpoint = "http://192.168.91.87:3000"

	data, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Secret-Key", os.Getenv("BOT_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}
