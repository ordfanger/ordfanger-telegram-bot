package service

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type Record struct {
	Word         string   `json:"word"`
	Language     string   `json:"language"`
	PartOfSpeech string   `json:"part_of_speech"`
	Sentences    []string `json:"sentences"`
}

var logger = logrus.New()

func RecordNewWord(word string) {
	logger.Formatter = &logrus.JSONFormatter{}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	record := &Record{
		Word:         word,
		Language:     "en",
		PartOfSpeech: "noun",
		Sentences:    []string{"Hello World ", "Hi world"},
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		logger.Errorf("Got error marshalling map: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("WORDS_TABLE")),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		logger.Errorf("Got error %s", err.Error())
	}
}
