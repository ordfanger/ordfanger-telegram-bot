package chat

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	assertion "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ChatMock struct {
	mock.Mock
}

func (chatMock *ChatMock) Send(chat *Chat, response *Responses) {
	chatMock.Called(chat, response)
}

func (chatMock *ChatMock) GetState(chat *Chat) error {
	args := chatMock.Called(chat)
	return args.Error(0)

}
func (chatMock *ChatMock) SaveState(chat *Chat) error {
	args := chatMock.Called(chat)
	return args.Error(0)
}
func (chatMock *ChatMock) RecordNewWord(chat *Chat) {
	chatMock.Called(chat)
}

func TestChatSend(t *testing.T) {
	chatMock := &ChatMock{}

	response := &Responses{}

	actual := &Chat{
		Bot: chatMock,
	}

	chatMock.On("Send", actual, response)

	actual.Send(response)

	chatMock.AssertExpectations(t)
	chatMock.AssertNumberOfCalls(t, "Send", 1)
}

func TestChatDecisionTree(t *testing.T) {
	var (
		MockReceivedWord                 = "MockReceivedWord"
		MockReceivedLanguage             = "MockReceivedLanguage"
		MockReceivedNonSupportedLanguage = "MockReceivedNonSupportedLanguage"
		MockReceivedPartOfSpeech         = "MockReceivedPartOfSpeech"
		MockReceivedUnknownPartOfSpeech  = "MockReceivedUnknownPartOfSpeech"
		MockReceivedSentences            = "MockReceivedSentences"
	)

	tt := []struct {
		title    string
		chat     *Chat
		response *Responses
		mockType string
	}{
		{
			"test isCommand: (start)",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "/start",
						Entities: &[]tgbotapi.MessageEntity{
							{Type: "bot_command", Offset: 0, Length: 6},
						},
					},
				},
			},
			&Responses{
				Text: "Hey! Seems you have new words.\nLet's save it!",
			},
			"",
		},
		{
			"test isCommand: (unknown)",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "/unknown",
						Entities: &[]tgbotapi.MessageEntity{
							{Type: "bot_command", Offset: 0, Length: 8},
						},
					},
				},
			},
			&Responses{
				Text: "Unknown command.",
			},
			"",
		},
		{
			"test received word",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "test",
					},
				},
				State: &State{
					Step: 1,
				},
			},
			&Responses{
				Text:                "Got you word. Pick language to save this word.",
				ReplyKeyboardMarkup: LanguageKeyboard(),
			},
			MockReceivedWord,
		},
		{
			"test received language",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "EN",
					},
				},
				State: &State{
					Step: 2,
				},
			},
			&Responses{
				Text:                "Nice, then choose part of speech.",
				ReplyKeyboardMarkup: PartOfSpeech(EN),
			},
			MockReceivedLanguage,
		},
		{
			"test received not supported language",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "SOME",
					},
				},
				State: &State{
					Step: 2,
				},
			},
			&Responses{
				Text:                "Your language is not supported. Please, select from keyboard.",
				ReplyKeyboardMarkup: LanguageKeyboard(),
			},
			MockReceivedNonSupportedLanguage,
		},
		{
			"test received part of speech",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "noun",
					},
				},
				State: &State{
					Step: 3,
					UserInputs: Record{
						Language: "EN",
					},
				},
			},
			&Responses{
				Text: "Type sentences to see how to use your word.",
			},
			MockReceivedPartOfSpeech,
		},
		{
			"test received unknown part of speech",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "test",
					},
				},
				State: &State{
					Step: 3,
					UserInputs: Record{
						Language: "EN",
					},
				},
			},
			&Responses{
				Text:                "Please, select part of speech.",
				ReplyKeyboardMarkup: PartOfSpeech(EN),
			},
			MockReceivedUnknownPartOfSpeech,
		},
		{
			"test received sentences",
			&Chat{
				Update: &Update{
					&tgbotapi.Message{
						Text: "hello world\nHello world!",
					},
				},
				State: &State{
					Step: 4,
					UserInputs: Record{
						Language: "EN",
					},
				},
			},
			&Responses{
				Text: "Brilliant! All info has been saved.\nIf you need to save another word, just type it.",
			},
			MockReceivedSentences,
		},
	}

	for _, tc := range tt {
		t.Run(tc.title, func(t *testing.T) {
			assert := assertion.New(t)

			chatMock := &ChatMock{}

			switch tc.mockType {
			case MockReceivedWord:
				chatMock.On("SaveState", tc.chat).Return(nil)

				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(2, tc.chat.State.Step)
				assert.Equal("test", tc.chat.State.UserInputs.Word)
				assert.Equal(tc.response, response)

				chatMock.AssertExpectations(t)
				chatMock.AssertNumberOfCalls(t, "SaveState", 1)

			case MockReceivedLanguage:
				chatMock.On("SaveState", tc.chat).Return(nil)

				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(3, tc.chat.State.Step)
				assert.Equal("EN", tc.chat.State.UserInputs.Language)
				assert.Equal(tc.response, response)

				chatMock.AssertExpectations(t)
				chatMock.AssertNumberOfCalls(t, "SaveState", 1)

			case MockReceivedNonSupportedLanguage:
				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(2, tc.chat.State.Step)
				assert.Equal(tc.response, response)

				chatMock.AssertNotCalled(t, "SaveState")

			case MockReceivedPartOfSpeech:
				chatMock.On("SaveState", tc.chat).Return(nil)
				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(tc.response, response)
				assert.Equal("noun", tc.chat.State.UserInputs.PartOfSpeech)
				assert.Equal(4, tc.chat.State.Step)

				chatMock.AssertNumberOfCalls(t, "SaveState", 1)
				chatMock.AssertExpectations(t)

			case MockReceivedUnknownPartOfSpeech:
				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(3, tc.chat.State.Step)
				assert.Equal(tc.response, response)

				chatMock.AssertNotCalled(t, "SaveState")

			case MockReceivedSentences:
				chatMock.On("SaveState", tc.chat).Return(nil)
				chatMock.On("RecordNewWord", tc.chat)

				tc.chat.Connection = chatMock

				response := tc.chat.DecisionTree()

				assert.Equal(1, tc.chat.State.Step)
				assert.Equal(tc.response, response)
				assert.Equal([]string{"hello world", "Hello world!"}, tc.chat.State.UserInputs.Sentences)

				chatMock.AssertNumberOfCalls(t, "RecordNewWord", 1)
				chatMock.AssertNumberOfCalls(t, "SaveState", 1)
			default:
				response := tc.chat.DecisionTree()
				assert.Equal(tc.response, response)
			}

		})
	}
}
