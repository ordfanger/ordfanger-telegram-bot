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
func DecisionTree(connection *dynamodb.DynamoDB, state *State) error {
	if state.Step == ReceivedWord {
		// create bot response
		// send bot response
		// move to next step

		state.Step = 2

		err := SaveState(connection, state)
		logger.Error(err)

		return nil
	}

	if state.Step == ReceivedLanguage {
		state.Step = 3
		err := SaveState(connection, state)
		logger.Error(err)

		return nil
	}

	if state.Step == ReceivedPartOfSpeech {
		state.Step = 4
		err := SaveState(connection, state)
		logger.Error(err)

		return nil
	}

	if state.Step == ReceivedSentences {
		state.Step = 1
		err := SaveState(connection, state)
		logger.Error(err)
	}

	return nil
}
