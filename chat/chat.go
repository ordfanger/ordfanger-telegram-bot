package chat

import (
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

const (
	ReceivedWord = iota + 1
	ReceivedLanguage
	ReceivedPartOfSpeech
	ReceivedSentences
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

type Chat struct {
	Logger     *logrus.Logger
	Bot        *tgbotapi.BotAPI
	Connection *dynamodb.DynamoDB
	Update     *Update
	State      *State
}

type Record struct {
	ID           string   `json:"id"`
	Word         string   `json:"word"`
	Language     string   `json:"language"`
	PartOfSpeech string   `json:"part_of_speech"`
	Sentences    []string `json:"sentences"`
}

func (chat *Chat) DecisionTree() (*Responses, error) {
	if chat.Update.Message.IsCommand() {
		command := chat.Update.Message.Command()
		response := &Responses{}

		if command == "start" {
			response.Text = string(Welcome)
			return response, nil
		}

		response.Text = string(UnknownCommand)
		return response, nil
	}

	if chat.State.Step == ReceivedWord {
		chat.State.Step = 2
		chat.State.UserInputs.Word = chat.Update.Message.Text

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text:                string(OnReceivedWord),
			ReplyKeyboardMarkup: LanguageKeyboard(),
		}, nil
	}

	if chat.State.Step == ReceivedLanguage {
		text := chat.Update.Message.Text

		language, err := GetLanguageFromText(text)
		if err != nil {
			return &Responses{
				Text:                string(UnknownLanguage),
				ReplyKeyboardMarkup: LanguageKeyboard(),
			}, nil
		}

		chat.State.Step = 3
		chat.State.UserInputs.Language = text

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text:                string(OnReceivedLanguage),
			ReplyKeyboardMarkup: PartOfSpeech(language),
		}, nil
	}

	if chat.State.Step == ReceivedPartOfSpeech {
		text := chat.Update.Message.Text
		language, _ := GetLanguageFromText(chat.State.UserInputs.Language)

		if !CheckIfPartOfSpeechExists(language, text) {
			return &Responses{
				Text:                string(UnknownPartOfSpeech),
				ReplyKeyboardMarkup: PartOfSpeech(language),
			}, nil
		}

		chat.State.Step = 4
		chat.State.UserInputs.PartOfSpeech = text

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text: string(OnReceivedPartOfSpeech),
		}, nil
	}

	if chat.State.Step == ReceivedSentences {
		chat.State.Step = 1

		sentences := strings.Split(chat.Update.Message.Text, "\n")
		chat.State.UserInputs.Sentences = sentences

		chat.RecordNewWord()

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state %v", err)
		}
	}

	return &Responses{
		Text: string(OnReceivedSentences),
	}, nil
}

func (chat *Chat) Send(response *Responses) {
	msg := tgbotapi.NewMessage(chat.State.ChatID, response.Text)

	if response.ReplyKeyboardMarkup != nil {
		msg.ReplyMarkup = response.ReplyKeyboardMarkup
	}

	if response.ReplyKeyboardMarkup == nil {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	if _, err := chat.Bot.Send(msg); err != nil {
		chat.Logger.Error(err)
	}
}

func (chat *Chat) GetState() error {
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

	out, err := chat.Connection.Query(params)
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

func (chat *Chat) SaveState() error {
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

	_, err = chat.Connection.PutItem(input)
	if err != nil {
		chat.Logger.Errorf("error while saving state: %v", err.Error())
		return err
	}

	return nil
}

func (chat *Chat) RecordNewWord() {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	uuId := uuid.NewV4()

	record := &Record{
		ID:           uuId.String(),
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
