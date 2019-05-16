package chat

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ReceivedWord = iota + 1
	ReceivedLanguage
	ReceivedPartOfSpeech
	ReceivedSentences
)

func DecisionTree(connection *dynamodb.DynamoDB, state *State, message *tgbotapi.Message) (*Responses, error) {
	if message.IsCommand() {
		command := message.Command()
		response := &Responses{}

		if command == "start" {
			response.Text = string(Welcome)
			return response, nil
		}

		response.Text = string(UnknownCommand)
		return response, nil
	}

	if state.Step == ReceivedWord {
		state.Step = 2
		state.UserInputs.Word = message.Text

		if err := SaveState(connection, state); err != nil {
			logger.Errorf("can't save state", err)
		}

		return &Responses{
			Text:                string(OnReceivedWord),
			ReplyKeyboardMarkup: LanguageKeyboard(),
		}, nil
	}

	if state.Step == ReceivedLanguage {
		text := message.Text

		language, err := GetLanguageFromText(text)
		if err != nil {
			return &Responses{
				Text:                string(UnknownLanguage),
				ReplyKeyboardMarkup: LanguageKeyboard(),
			}, nil
		}

		state.Step = 3
		state.UserInputs.Language = text

		if err := SaveState(connection, state); err != nil {
			logger.Errorf("can't save state", err)
		}

		return &Responses{
			Text:                string(OnReceivedLanguage),
			ReplyKeyboardMarkup: PartOfSpeech(language),
		}, nil
	}

	if state.Step == ReceivedPartOfSpeech {
		text := message.Text
		language, _ := GetLanguageFromText(state.UserInputs.Language)

		if !CheckIfPartOfSpeechExists(language, text) {
			return &Responses{
				Text:                string(UnknownPartOfSpeech),
				ReplyKeyboardMarkup: PartOfSpeech(language),
			}, nil
		}

		state.Step = 4
		state.UserInputs.PartOfSpeech = text

		if err := SaveState(connection, state); err != nil {
			logger.Errorf("can't save state", err)
		}

		return &Responses{
			Text: string(OnReceivedPartOfSpeech),
		}, nil
	}

	if state.Step == ReceivedSentences {
		state.Step = 1

		sentences := strings.Split(message.Text, "\n")
		state.UserInputs.Sentences = sentences

		RecordNewWord(state)

		if err := SaveState(connection, state); err != nil {
			logger.Errorf("can't save state", err)
		}
	}

	return &Responses{
		Text: string(OnReceivedSentences),
	}, nil
}
