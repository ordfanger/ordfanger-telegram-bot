package chat

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

func GetChatState(context *Context) (*State, error) {
	chatState := &State{}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("CHAT_STATE_TABLE")),
		KeyConditionExpression: aws.String("userID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				N: aws.String(strconv.Itoa(context.Update.Message.From.ID)),
			},
		},
	}

	out, err := context.Connection.Query(params)
	if err != nil {
		context.Logger.Errorf("querying state failed: %v", err.Error())
		return nil, err
	}

	if len(out.Items) == 0 {
		return &State{
			Step:          1,
			UserID:        context.Update.Message.From.ID,
			UserFirstName: context.Update.Message.From.FirstName,
			UserLastName:  context.Update.Message.From.LastName,
			UserName:      context.Update.Message.From.UserName,
			ChatID:        context.Update.Message.Chat.ID,
			UserInputs:    Record{},
		}, nil
	}

	for _, item := range out.Items {
		if err := dynamodbattribute.UnmarshalMap(item, chatState); err != nil {
			context.Logger.Errorf("unmarshalMap failed: %v", err.Error())
			return nil, err
		}
	}

	return chatState, nil
}

func SaveState(context *Context) error {
	context.Logger.WithFields(logrus.Fields{"state": context.State}).Info("saving state")

	av, err := dynamodbattribute.MarshalMap(context.State)
	if err != nil {
		context.Logger.Errorf("error marshalling state: %v", err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("CHAT_STATE_TABLE")),
	}

	_, err = context.Connection.PutItem(input)
	if err != nil {
		context.Logger.Errorf("error while saving state: %v", err.Error())
		return err
	}

	return nil
}
