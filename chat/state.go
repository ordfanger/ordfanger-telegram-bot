package chat

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ordfanger/ordfanger-telegram-bot/internal"
	"github.com/sirupsen/logrus"
)

var logger = internal.NewLogger()

type State struct {
	Step          int    `json:"step"`
	UserID        int    `json:"userID"`
	ChatID        int64  `json:"chatID"`
	UserFirstName string `json:"first_name"`
	UserLastName  string `json:"last_name"`
	UserName      string `json:"username"`
	UserInputs    Record `json:"user_inputs"`
}

func GetChatState(connection *dynamodb.DynamoDB, message *tgbotapi.Message) (*State, error) {
	chatState := &State{}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("CHAT_STATE_TABLE")),
		KeyConditionExpression: aws.String("userID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				N: aws.String(strconv.Itoa(message.From.ID)),
			},
		},
	}

	out, err := connection.Query(params)
	if err != nil {
		logger.Errorf("querying state failed: %v", err.Error())
		return nil, err
	}

	if len(out.Items) == 0 {
		return &State{
			Step:          1,
			UserID:        message.From.ID,
			UserFirstName: message.From.FirstName,
			UserLastName:  message.From.LastName,
			UserName:      message.From.UserName,
			ChatID:        message.Chat.ID,
			UserInputs:    Record{},
		}, nil
	}

	for _, item := range out.Items {
		if err := dynamodbattribute.UnmarshalMap(item, chatState); err != nil {
			logger.Errorf("unmarshalMap failed: %v", err.Error())
			return nil, err
		}
	}

	return chatState, nil
}

func SaveState(connection *dynamodb.DynamoDB, state *State) error {
	logger.WithFields(logrus.Fields{"state": state}).Info("saving state")

	av, err := dynamodbattribute.MarshalMap(state)
	if err != nil {
		logger.Errorf("error marshalling state: %v", err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("CHAT_STATE_TABLE")),
	}

	_, err = connection.PutItem(input)
	if err != nil {
		logger.Errorf("error while saving state: %v", err.Error())
		return err
	}

	return nil
}
