package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Request用オブジェクト
type TechplayParam struct {
	EventId int `json:"event_id"`
}

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
	Total      string  `json:"total"`
	DetailList Details `json:"detail_list"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(request.Body)

	// イベントURL
	eventUrl := "https://techplay.jp/event/%d"

	// リクエスト情報をパース
	requestParam := new(TechplayParam)
	if err := json.Unmarshal(([]byte)(request.Body), requestParam); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// イベントページ取得
	doc, err := goquery.NewDocument(fmt.Sprintf(eventUrl, requestParam.EventId))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// TechplayEvent生成
	techplayEvent := TechplayEvent{
		EventUrl:   fmt.Sprintf(eventUrl, requestParam.EventId),
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

	// 合計人数取得
	techplayEvent.Total = getTotal(techplayEvent.DetailList)

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
	result = strings.Replace(result, "／", "/", -1)

	// connpass-apiと同じ表記にするために最初の"人"を除外
	if idx := strings.Index(result, "/"); 0 < idx {
		result = strings.Replace(result[:idx], "人", "", 1) + result[idx:]
	}

	return result
}

// 合計人数を取得する
func getTotal(dArr Details) string {
	total := 0
	for _, d := range dArr {
		if num, result := capacity2JoinNum(d.Capacity); result {
			log.Printf("%s, %d¥n", d.Capacity, num)
			total += num
		}
	}
	return fmt.Sprintf("%d人", total)
}

// Detail.Capacityから参加人数のみ取得
func capacity2JoinNum(target string) (int, bool) {
	if idx := strings.Index(target, "/"); 0 < idx {
		mol := 0
		if num, err := strconv.Atoi(target[0:idx]); err == nil {
			mol = num
		}
		if idxDen := strings.Index(target, "人"); 0 < idxDen {
			if num, err := strconv.Atoi(target[idx+1 : idxDen]); err == nil && num < mol {
				mol = num
			}
		}
		return mol, true
	} else if idx := strings.Index(target, "人"); 0 < idx {
		if num, err := strconv.Atoi(target[0:idx]); err == nil {
			return num, true
		}
	}
	return -1, false
}
