package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(request.Body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "slack-app-sample",
	}, nil
}

func main() {
	lambda.Start(handler)
}
