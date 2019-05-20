package chat

import (
	"strings"
)

const (
	ReceivedWord = iota + 1
	ReceivedLanguage
	ReceivedPartOfSpeech
	ReceivedSentences
)

func DecisionTree(context *Context) (*Responses, error) {
	if context.Update.Message.IsCommand() {
		command := context.Update.Message.Command()
		response := &Responses{}

		if command == "start" {
			response.Text = string(Welcome)
			return response, nil
		}

		response.Text = string(UnknownCommand)
		return response, nil
	}

	if context.State.Step == ReceivedWord {
		context.State.Step = 2
		context.State.UserInputs.Word = context.Update.Message.Text

		if err := SaveState(context); err != nil {
			context.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text:                string(OnReceivedWord),
			ReplyKeyboardMarkup: LanguageKeyboard(),
		}, nil
	}

	if context.State.Step == ReceivedLanguage {
		text := context.Update.Message.Text

		language, err := GetLanguageFromText(text)
		if err != nil {
			return &Responses{
				Text:                string(UnknownLanguage),
				ReplyKeyboardMarkup: LanguageKeyboard(),
			}, nil
		}

		context.State.Step = 3
		context.State.UserInputs.Language = text

		if err := SaveState(context); err != nil {
			context.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text:                string(OnReceivedLanguage),
			ReplyKeyboardMarkup: PartOfSpeech(language),
		}, nil
	}

	if context.State.Step == ReceivedPartOfSpeech {
		text := context.Update.Message.Text
		language, _ := GetLanguageFromText(context.State.UserInputs.Language)

		if !CheckIfPartOfSpeechExists(language, text) {
			return &Responses{
				Text:                string(UnknownPartOfSpeech),
				ReplyKeyboardMarkup: PartOfSpeech(language),
			}, nil
		}

		context.State.Step = 4
		context.State.UserInputs.PartOfSpeech = text

		if err := SaveState(context); err != nil {
			context.Logger.Errorf("can't save state  %v", err)
		}

		return &Responses{
			Text: string(OnReceivedPartOfSpeech),
		}, nil
	}

	if context.State.Step == ReceivedSentences {
		context.State.Step = 1

		sentences := strings.Split(context.Update.Message.Text, "\n")
		context.State.UserInputs.Sentences = sentences

		RecordNewWord(context)

		if err := SaveState(context); err != nil {
			context.Logger.Errorf("can't save state %v", err)
		}
	}

	return &Responses{
		Text: string(OnReceivedSentences),
	}, nil
}
