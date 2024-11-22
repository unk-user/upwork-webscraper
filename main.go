package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const baseUrl = "https://www.upwork.com/nx/search/jobs/"

type Response struct {
	StatusCode int   `json:"statusCode"`
	Length     int   `json:"length"`
	Jobs       []Job `json:"jobs"`
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

	log.Println("Done")
	return Response{StatusCode: 200, Jobs: jobs, Length: len(jobs)}, nil
}

func main() {
	lambda.Start(handler)
}
