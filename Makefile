build:
	mkdir -p functions
	go get ./source/connpass-api
	go build -o ./functions/connpass-api ./source/connpass-api
	go get ./source/techplay
	go build -o ./functions/techplay ./source/techplay
	go get ./source/slack-app-sample
	go build -o ./functions/slack-app-sample ./source/slack-app-sample