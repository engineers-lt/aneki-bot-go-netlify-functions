package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Detail struct {
	Category string `json:category`
	Capacity string `json:capacity`
}
type Details []*Detail
type TechplayEvent struct {
	Title      string  `json:title`
	Day        string  `json:day`
	Time       string  `json:time`
	DetailList Details `json:detail_list`
}

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

	// TechplayEvent生成
	techplayEvent := TechplayEvent{DetailList: Details{}}

	// タイトル部分抜き出し
	//  root: <div class="title-heading">
	techplayEvent.Title = doc.Find("div.title-heading > h1").Text()

	// 参加枠、定員抜き出し
	//  root: <div id="participationTable">
	tableSelection := doc.Find("div#participationTable > table > tbody > tr")
	reSpan := regexp.MustCompile(`(?m)<.span.>`)
	//reNewLine := regexp.MustCompile(`(?m)\n`)
	tableSelection.Each(func(_ int, s *goquery.Selection) {
		d := Detail{Category: s.Find("td.category > div.category-inner > div").Text()}
		c := reSpan.ReplaceAllString(strings.Replace(s.Find("td.capacity").Text(), "\n", "", -1), "")
		d.Capacity = strings.Replace(strings.Replace(c, " ", "", -1), "定員", "", -1)
		techplayEvent.DetailList = append(techplayEvent.DetailList, &d)
	})

	// jsonエンコード
	result, _ := json.Marshal(techplayEvent)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(result),
	}, nil
}

func main() {
	lambda.Start(handler)
}
