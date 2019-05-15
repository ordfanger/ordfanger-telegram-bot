package chat

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	ReceivedWord = iota + 1
	ReceivedLanguage
	ReceivedPartOfSpeech
	ReceivedSentences
)

// try to create context and set connection in context.
func DecisionTree(connection *dynamodb.DynamoDB, state *State) (*Responses, error) {
	if state.Step == ReceivedWord {
		// create bot response
		// send bot response
		// move to next step

		state.Step = 2

		_ = SaveState(connection, state)

		return &Responses{Text: "Got you word", ReplyKeyboardMarkup: LanguageKeyboard()}, nil
	}

	if state.Step == ReceivedLanguage {
		state.Step = 3
		_ = SaveState(connection, state)

		return &Responses{
			Text:                "Got you language, select part of speech",
			ReplyKeyboardMarkup: PartOfSpeech(Language(EN)),
		}, nil
	}

	if state.Step == ReceivedPartOfSpeech {
		state.Step = 4
		_ = SaveState(connection, state)

		return &Responses{
			Text: "Cool, enter sentences",
		}, nil
	}

	if state.Step == ReceivedSentences {
		state.Step = 1
		_ = SaveState(connection, state)
	}

	return &Responses{
		Text: "Done!",
	}, nil
}
