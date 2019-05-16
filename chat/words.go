package chat

import (
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Record struct {
	ID           string   `json:"id"`
	Word         string   `json:"word"`
	Language     string   `json:"language"`
	PartOfSpeech string   `json:"part_of_speech"`
	Sentences    []string `json:"sentences"`
}

func RecordNewWord(state *State) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	uuId := uuid.NewV4()

	record := &Record{
		ID:           uuId.String(),
		Word:         state.UserInputs.Word,
		Language:     state.UserInputs.Language,
		PartOfSpeech: state.UserInputs.PartOfSpeech,
		Sentences:    state.UserInputs.Sentences,
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		logger.Errorf("error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("WORDS_TABLE")),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		logger.Errorf("can't save the word %s", err.Error())
	}
}
