package main

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// connpassイベントサーチAPiを実行した際に返却されるJSONオブジェクト
type ConnpassApiObject struct {
	ResultsReturned int `json:"results_returned"`
	Events          []struct {
		EventURL      string `json:"event_url"`
		EventType     string `json:"event_type"`
		OwnerNickname string `json:"owner_nickname"`
		Series        struct {
			URL   string `json:"url"`
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"series"`
		UpdatedAt        time.Time `json:"updated_at"`
		Lat              string    `json:"lat"`
		StartedAt        time.Time `json:"started_at"`
		HashTag          string    `json:"hash_tag"`
		Title            string    `json:"title"`
		EventID          int       `json:"event_id"`
		Lon              string    `json:"lon"`
		Waiting          int       `json:"waiting"`
		Limit            int       `json:"limit"`
		OwnerID          int       `json:"owner_id"`
		OwnerDisplayName string    `json:"owner_display_name"`
		Description      string    `json:"description"`
		Address          string    `json:"address"`
		Catch            string    `json:"catch"`
		Accepted         int       `json:"accepted"`
		EndedAt          time.Time `json:"ended_at"`
		Place            string    `json:"place"`
	} `json:"events"`
	ResultsStart     int `json:"results_start"`
	ResultsAvailable int `json:"results_available"`
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

	// イベントURL
	eventUrl := "https://engineers.connpass.com/event/104940/"

	// イベントページ取得
	// connpassApiObj := getConnpassApiObject(eventUrl)
	doc, err := goquery.NewDocument(eventUrl)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// ConnpassEvent生成
	connpassEvent := ConnpassEvent{
		EventUrl:   eventUrl,
		DetailList: Details{},
	}

	// TODO: タイトル部分抜き出し
	// どうにかせねば(´・ω・｀)
	connpassEvent.Title = "秋の夜長の自由研究LT大会"
	connpassEvent.EventUrl = eventUrl

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
