package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// リクエスト情報をログ出力
	log.Print(request.Body)

	// 疎通確認用のリクエストか確認
	slackRequestFirst := new(SlackRequestFirst)
	if err := json.Unmarshal(([]byte)(request.Body), slackRequestFirst); err == nil && 0 < len(slackRequestFirst.Challenge) {
		log.Print("first")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       slackRequestFirst.Challenge,
		}, nil
	} else {
		slackRequest := new(SlackRequest)
		if err := json.Unmarshal(([]byte)(request.Body), slackRequest); err == nil && 0 < len(slackRequest.Token) {
			// トークンが取得できたらオウム返しを行う
			token := os.Getenv("BOT_TOKEN")
			api := slack.New(token)
			_, _, err = api.PostMessage(
				slackRequest.Event.Channel,
				slack.MsgOptionText(
					slackRequest.Event.Text,
					false,
				),
			)

			// 返す先はないが200OKとする
			if err != nil {
				log.Print("err: ", err)
			} else {
				log.Print("ok")
			}
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "",
			}, nil
		} else {
			log.Print("err")
			return events.APIGatewayProxyResponse{}, nil
		}
	}
}

func main() {
	lambda.Start(handler)
}
