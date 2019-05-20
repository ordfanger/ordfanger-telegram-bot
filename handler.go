package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ordfanger/ordfanger-telegram-bot/server"
)

// lambda main handler
func main() {
	lambda.Start(server.Server)
}
