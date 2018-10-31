build:
	mkdir -p functions
	go get ./source/connpass-api
	go build -o ./functions/connpass-api ./source/connpass-api
	go get ./source/slack-api-sample
	go build -o ./functions/slack-app-sample ./source/slack-app-sample