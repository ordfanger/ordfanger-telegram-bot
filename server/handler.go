package main

import (
	"context"
	"encoding/json"
	"github.com/ordfanger/ordfanger-telegram-bot/service"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	logger "github.com/sirupsen/logrus"
)

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

type Response events.APIGatewayProxyResponse

type Update struct {
	Message *tgbotapi.Message
}

// Server add comment here
func Server(_ context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	botAPIKey := os.Getenv("BOT_API_KEY")

	var update Update

	decoder := json.NewDecoder(strings.NewReader(req.Body))

	if err := decoder.Decode(&update); err != nil {
		logger.Error(err)
	}

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		logger.Error(err)
	}

	bot.Debug = true

	logger.Infof("Authorized on account %s", bot.Self.UserName)

	logger.WithFields(logger.Fields{
		"update": &update,
	}).Info("New update")

	service.RecordNewWord(update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		logger.Error(err)
	}

	return Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Server)
}
