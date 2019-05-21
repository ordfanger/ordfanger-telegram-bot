package chat

import (
	"strings"

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
	Bot        BotAPI
	Connection PersistenceLayer
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

func (chat *Chat) DecisionTree() *Responses {
	if chat.Update.Message.IsCommand() {
		command := chat.Update.Message.Command()
		response := &Responses{}

		if command == "start" {
			response.Text = string(Welcome)
			return response
		}

		response.Text = string(UnknownCommand)
		return response
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
		}
	}

	if chat.State.Step == ReceivedLanguage {
		text := chat.Update.Message.Text

		language, err := GetLanguageFromText(text)
		if err != nil {
			return &Responses{
				Text:                string(UnknownLanguage),
				ReplyKeyboardMarkup: LanguageKeyboard(),
			}
		}

		chat.State.Step = 3
		chat.State.UserInputs.Language = text

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text:                string(OnReceivedLanguage),
			ReplyKeyboardMarkup: PartOfSpeech(language),
		}
	}

	if chat.State.Step == ReceivedPartOfSpeech {
		text := chat.Update.Message.Text
		language, _ := GetLanguageFromText(chat.State.UserInputs.Language)

		if !CheckIfPartOfSpeechExists(language, text) {
			return &Responses{
				Text:                string(UnknownPartOfSpeech),
				ReplyKeyboardMarkup: PartOfSpeech(language),
			}
		}

		chat.State.Step = 4
		chat.State.UserInputs.PartOfSpeech = text

		if err := chat.SaveState(); err != nil {
			chat.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text: string(OnReceivedPartOfSpeech),
		}
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
	}
}

func (chat *Chat) Send(response *Responses) {
	chat.Bot.Send(chat, response)
}

func (chat *Chat) GetState() error {
	err := chat.Connection.GetState(chat)
	return err
}

func (chat *Chat) SaveState() error {
	err := chat.Connection.SaveState(chat)
	return err
}

func (chat *Chat) RecordNewWord() {
	chat.Connection.RecordNewWord(chat)
}
