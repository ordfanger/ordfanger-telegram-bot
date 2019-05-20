package chat

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Update struct {
	Message *tgbotapi.Message
}

type State struct {
	Step          int    `json:"step"`
	UserID        int    `json:"userID"`
	ChatID        int64  `json:"chatID"`
	UserFirstName string `json:"first_name"`
	UserLastName  string `json:"last_name"`
	UserName      string `json:"username"`
	UserInputs    Record `json:"user_inputs"`
}

type Context struct {
	Logger     *logrus.Logger
	Bot        *tgbotapi.BotAPI
	Connection *dynamodb.DynamoDB
	Update     *Update
	State      *State
}
