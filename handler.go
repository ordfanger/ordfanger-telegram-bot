package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ordfanger/ordfanger-telegram-bot/server"
)

func main() {
	lambda.Start(server.Server)
}
