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

type SlackRequest struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		Type        string `json:"type"`
		User        string `json:"user"`
		Text        string `json:"text"`
		ClientMsgID string `json:"client_msg_id"`
		Ts          string `json:"ts"`
		Channel     string `json:"channel"`
		EventTs     string `json:"event_ts"`
	} `json:"event"`
	Type        string   `json:"type"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
}

type SlackResponse struct {
	Token   string `json:"text"`
	Channel string `json:"text"`
	Text    string `json:"text"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// リクエスト情報をログ出力
	log.Print(request.Body)

	// 疎通確認用のリクエストか確認
	slackRequestFirst := new(SlackRequestFirst)
	if err := json.Unmarshal(([]byte)(request.Body), slackRequestFirst); err == nil && 0 < len(slackRequestFirst.Challenge) {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       slackRequestFirst.Challenge,
		}, nil
	} else {
		slackRequest := new(SlackRequest)
		if err := json.Unmarshal(([]byte)(request.Body), slackRequest); err == nil && 0 < len(slackRequest.Token) {
			slackResponse := SlackResponse{
				slackRequest.Token,
				slackRequest.Event.Channel,
				"受信した。",
			}
			result, _ := json.Marshal(slackResponse)
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       string(result),
			}, nil
		} else {
			return events.APIGatewayProxyResponse{}, nil
		}
	}
}

func main() {
	lambda.Start(handler)
}
