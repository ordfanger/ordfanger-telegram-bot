package server

import (
	"context"
	"encoding/json"

	"github.com/ordfanger/ordfanger-telegram-bot/chat"
	"github.com/ordfanger/ordfanger-telegram-bot/internal"

	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

// Response for match APIGateway Proxy Response.
type Response events.APIGatewayProxyResponse

// Starting poing for chat. All needed instanses created here.
func Server(_ context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	botAPIKey := os.Getenv("BOT_API_KEY")

	chatContext := &chat.Chat{
		Logger: internal.NewLogger(),
		Connection: &chat.DB{
			DynamoDB: internal.NewDBConnection(),
		},
	}

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.Bot = &chat.Bot{
		BotAPI: bot,
	}

	decoder := json.NewDecoder(strings.NewReader(req.Body))
	if err := decoder.Decode(&chatContext.Update); err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.Logger.Infof("authorized on account %s", bot.Self.UserName)
	chatContext.Logger.WithFields(logrus.Fields{"update": &chatContext.Update}).Info("received a new update")

	err = chatContext.GetState()
	if err != nil {
		chatContext.Logger.Error(err)
	}

	response := chatContext.DecisionTree()
	if err != nil {
		chatContext.Logger.Error(err)
	}

	chatContext.Send(response)

	return Response{StatusCode: 200}, nil
}
