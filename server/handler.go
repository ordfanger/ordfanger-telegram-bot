package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ordfanger/ordfanger-telegram-bot/chat"
	"github.com/ordfanger/ordfanger-telegram-bot/internal"

	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var logger = internal.NewLogger()

type Response events.APIGatewayProxyResponse

type Update struct {
	Message *tgbotapi.Message
}

func NewDBConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	return svc
}

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

	logger.Infof("authorized on account %s", bot.Self.UserName)
	logger.WithFields(logrus.Fields{"update": &update}).Info("received a new update")

	connection := NewDBConnection()
	chatState, err := chat.GetChatState(connection, update.Message)
	if err != nil {
		logger.Error(err)
	}

	response, err := chat.DecisionTree(connection, chatState, update.Message)
	if err != nil {
		logger.Error(err)
	}

	chat.Send(bot, response, chatState)

	return Response{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Server)
}
