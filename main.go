package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	payload := Payload{
		CategoryId: event.QueryStringParameters["categoryId"],
		Keywords:   event.QueryStringParameters["keywords"],
	}

	jobs, err := GetNewJobs(payload)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	body, err := json.Marshal(jobs)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
