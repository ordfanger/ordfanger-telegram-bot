package service

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	logger "github.com/sirupsen/logrus"
)

type Record struct {
	Word         string   `json:"word"`
	Language     string   `json:"language"`
	PartOfSpeech string   `json:"part_of_speech"`
	Sentences    []string `json:"sentences"`
}

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
}

func RecordNewWord(word string) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	record := &Record{
		Word:         word,
		Language:     "en",
		PartOfSpeech: "noun",
		Sentences:    []string{"Hello World", "Hi world"},
	}

	logger.Info("GOT word")

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		logger.Errorf("Got error marshalling map: %s", err.Error())
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		logger.Errorf("Got error  %s", err.Error())
	}
}
