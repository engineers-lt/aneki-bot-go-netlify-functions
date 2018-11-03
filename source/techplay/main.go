package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// イベントURL
	eventUrl := "https://techplay.jp/event/%d"

	// とりあえず固定
	eventId := 705867

	// イベントページ取得
	doc, err := goquery.NewDocument(fmt.Sprintf(eventUrl, eventId))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// タイトル部分抜き出し
	//  from: <div class="title-heading">
	//  to: <div
	title, _ := doc.Find("div.title-heading > h1").Html()

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       title,
	}, nil
}

func main() {
	lambda.Start(handler)
}
