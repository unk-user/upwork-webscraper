package main

import (
	// "context"
	// "fmt"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-rod/rod/lib/launcher"
)


func main() {
	launcher := launcher.New()
	_, err := GetNewJobs("react", launcher)

	if err != nil {
		panic(err)
	}
}
