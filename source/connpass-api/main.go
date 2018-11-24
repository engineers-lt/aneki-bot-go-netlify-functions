package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Request用オブジェクト
type ConnpassParam struct {
	EventUrl string `json:"event_url"`
	Title    string `json:"title"`
}

// 返却するJSONオブジェクト
type Detail struct {
	Category string `json:"category"`
	Capacity string `json:"capacity"`
}
type Details []*Detail
type ConnpassEvent struct {
	EventUrl   string  `json:"event_url"`
	Title      string  `json:"title"`
	Day        string  `json:"day"`
	Time       string  `json:"time"`
	DetailList Details `json:"detail_list"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(request.Body)

	// リクエスト情報をパース
	requestParam := new(ConnpassParam)
	if err := json.Unmarshal(([]byte)(request.Body), requestParam); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// イベントページ取得
	doc, err := goquery.NewDocument(requestParam.EventUrl)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// ConnpassEvent生成
	// TODO: タイトルもスクレイピングで取得できるようにしたい
	connpassEvent := ConnpassEvent{
		EventUrl:   requestParam.EventUrl,
		Title:      requestParam.Title,
		DetailList: Details{},
	}

	// 日付部分抜き出し
	//  root: <div class="event-day">
	connpassEvent.Day = doc.Find("span.dtstart > p.ymd").Text()

	// 時間部分抜き出し
	//  root: <div class="event-time">
	connpassEvent.Time = textCleansing(nil, doc.Find("span.dtstart > span.hi").Text()+"〜"+doc.Find("span.dtend").Text())

	// 参加枠、定員抜き出し
	//  root: <div id="participationTable">
	tableSelection := doc.Find("div.ptypes > div.ptype")
	tableSelection.Each(func(_ int, s *goquery.Selection) {
		d := Detail{Category: s.Find("p.ptype_name").Text()}
		d.Capacity = textCleansing(nil, s.Find("p.participants > span.amount").Text()) + "人"
		connpassEvent.DetailList = append(connpassEvent.DetailList, &d)
	})

	// jsonエンコード
	result, _ := json.Marshal(connpassEvent)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(result),
	}, nil
}

func main() {
	lambda.Start(handler)
}

// テキストから改行およびスペースを削除
func textCleansing(r *regexp.Regexp, target string) string {
	result := strings.Replace(target, "\n", "", -1)
	result = strings.Replace(result, " ", "", -1)
	if r != nil {
		result = r.ReplaceAllString(result, "")
	}
	return strings.Replace(result, "／", "/", -1)
}
