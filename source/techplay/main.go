package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Detail struct {
	Category string `json:"category"`
	Capacity string `json:"capacity"`
}
type Details []*Detail
type TechplayEvent struct {
	EventUrl   string  `json:"event_url"`
	Title      string  `json:"title"`
	Day        string  `json:"day"`
	Time       string  `json:"time"`
	DetailList Details `json:"detail_list"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(request.Body)

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
	techplayEvent := TechplayEvent{
		EventUrl:   fmt.Sprintf(eventUrl, eventId),
		DetailList: Details{},
	}

	// タイトル部分抜き出し
	//  root: <div class="title-heading">
	techplayEvent.Title = doc.Find("div.title-heading > h1").Text()

	// 日付部分抜き出し
	//  root: <div class="event-day">
	techplayEvent.Day = doc.Find("div.event-day").Text()

	// 時間部分抜き出し
	//  root: <div class="event-time">
	techplayEvent.Time = doc.Find("div.event-time").Text()

	// 参加枠、定員抜き出し
	//  root: <div id="participationTable">
	tableSelection := doc.Find("div#participationTable > table > tbody > tr")
	reSpan := regexp.MustCompile(`<.span.>`)
	tableSelection.Each(func(_ int, s *goquery.Selection) {
		d := Detail{Category: s.Find("td.category > div.category-inner > div").Text()}
		d.Capacity = getCapacity(reSpan, s.Find("td.capacity").Text())
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

// データクレンジングをしながら定員のみの文字列を取得する
func getCapacity(r *regexp.Regexp, target string) string {
	result := strings.Replace(target, "\n", "", -1)
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "定員", "", -1)
	result = r.ReplaceAllString(result, "")
	return strings.Replace(result, "／", "/", -1)
}
