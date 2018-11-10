# SlackApp等で使用するAPIを定義

## ローカル環境でビルドする方法

```bash
GOOS=linux GOARCH=amd64 go build -o ./bin/connpass-api ./source/connpass-api
GOOS=linux GOARCH=amd64 go build -o ./bin/slack-app-sample ./source/slack-app-sample
GOOS=linux GOARCH=amd64 go build -o ./bin/techplay ./source/techplay
sam local start-api
```
