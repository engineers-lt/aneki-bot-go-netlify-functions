package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SlackRequestFirst struct {
	Token     string `json:token`
	Challenge string `json:challenge`
	Type      string `json:type`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// リクエスト情報をログ出力
	log.Print(request.Body)

	// jsonデコード
	slackRequest := new(SlackRequestFirst)
	jsonBytes := ([]byte)(request.Body)
	err := json.Unmarshal(jsonBytes, slackRequest)
	if err != nil {
		log.Fatal("Json Unmarshal error: ", err)

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       err.Error(),
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       slackRequest.Challenge,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}
