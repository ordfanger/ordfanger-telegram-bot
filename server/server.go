package server

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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

func newDBConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	return svc
}

func Server(_ context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	botAPIKey := os.Getenv("BOT_API_KEY")

	chatContext := &chat.Context{
		Logger:     internal.NewLogger(),
		Connection: newDBConnection(),
	}

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.Bot = bot

	decoder := json.NewDecoder(strings.NewReader(req.Body))
	if err := decoder.Decode(&chatContext.Update); err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.Logger.Infof("authorized on account %s", bot.Self.UserName)
	chatContext.Logger.WithFields(logrus.Fields{"update": &chatContext.Update}).Info("received a new update")

	chatState, err := chat.GetChatState(chatContext)
	if err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.State = chatState

	response, err := chat.DecisionTree(chatContext)
	if err != nil {
		chatContext.Logger.Error(err)
	}

	chat.Send(chatContext, response)

	return Response{StatusCode: 200}, nil
}
