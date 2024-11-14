package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func resp(body string, status int) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       body,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	const requiredQueryParameter = "keyword"

	keyword, ok := request.QueryStringParameters[requiredQueryParameter]
	if !ok {
		err := fmt.Sprintf(`Query parameter "%s" is required`, requiredQueryParameter)
		return resp(err, 400)
	}

	err := GetNewJobs(keyword)
	if err != nil {
		return resp(err.Error(), 500)
	}

	return resp("ok", 200)
}

func main() {
	lambda.Start(handler)
}
