package chat

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// PersistenceLayer interface for data storage.
type PersistenceLayer interface {
	GetState(chat *Chat) error
	SaveState(chat *Chat) error
	RecordNewWord(chat *Chat)
}

// DB is DynamoDB implamentation.
type DB struct {
	*dynamodb.DynamoDB
}

// GetState used to getting current chat state for bot decigions.
func (db *DB) GetState(chat *Chat) error {
	chatState := &State{}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("CHAT_STATE_TABLE")),
		KeyConditionExpression: aws.String("userID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				N: aws.String(strconv.Itoa(chat.Update.Message.From.ID)),
			},
		},
	}

	out, err := db.DynamoDB.Query(params)
	if err != nil {
		chat.Logger.Errorf("querying state failed: %v", err.Error())
		return err
	}

	if len(out.Items) == 0 {
		chat.State = &State{
			Step:          1,
			UserID:        chat.Update.Message.From.ID,
			UserFirstName: chat.Update.Message.From.FirstName,
			UserLastName:  chat.Update.Message.From.LastName,
			UserName:      chat.Update.Message.From.UserName,
			ChatID:        chat.Update.Message.Chat.ID,
			UserInputs:    Record{},
		}

		return nil
	}

	for _, item := range out.Items {
		if err := dynamodbattribute.UnmarshalMap(item, chatState); err != nil {
			chat.Logger.Errorf("unmarshalMap failed: %v", err.Error())
			return err
		}
	}

	chat.State = chatState

	return nil
}

// SaveState saves chat state.
func (db *DB) SaveState(chat *Chat) error {
	chat.Logger.WithFields(logrus.Fields{"state": chat.State}).Info("saving state")

	av, err := dynamodbattribute.MarshalMap(chat.State)
	if err != nil {
		chat.Logger.Errorf("error marshalling state: %v", err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("CHAT_STATE_TABLE")),
	}

	_, err = db.DynamoDB.PutItem(input)
	if err != nil {
		chat.Logger.Errorf("error while saving state: %v", err.Error())
		return err
	}

	return nil
}

// RecordNewWord stores final user's information about word.
func (db *DB) RecordNewWord(chat *Chat) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	uuID := uuid.NewV4()

	record := &Record{
		ID:           uuID.String(),
		Word:         chat.State.UserInputs.Word,
		Language:     chat.State.UserInputs.Language,
		PartOfSpeech: chat.State.UserInputs.PartOfSpeech,
		Sentences:    chat.State.UserInputs.Sentences,
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		chat.Logger.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("WORDS_TABLE")),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		chat.Logger.Errorf("can't save the word %s", err.Error())
	}
}
