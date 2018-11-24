# SlackApp等で使用するAPIを定義


# ローカル環境

## ビルド

```bash
GOOS=linux GOARCH=amd64 go build -o ./bin/connpass-api ./source/connpass-api
GOOS=linux GOARCH=amd64 go build -o ./bin/slack-app-sample ./source/slack-app-sample
GOOS=linux GOARCH=amd64 go build -o ./bin/techplay ./source/techplay
sam local start-api
```

## テスト

```bash
# techplay
curl http://127.0.0.1:3000/techplay -X POST -H "Content-Type: application/json" -d '{"event_id": "id"}'
# connpass-api
curl http://127.0.0.1:3000/connpass-api -X POST -H "Content-Type: application/json" -d '{"event_url": "url", "title": "title"}'
```
