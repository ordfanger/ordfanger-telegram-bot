package chat_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/ordfanger/ordfanger-telegram-bot/chat"
)

type BotAPIMock struct {
	mock.Mock
}

func (botMock *BotAPIMock) Send(chat *chat.Chat, response *chat.Responses) {
	botMock.Called(chat, response)
}

func TestChatSend(t *testing.T) {
	botMock := &BotAPIMock{}

	resp := &chat.Responses{}

	actual := &chat.Chat{
		Bot: botMock,
	}

	botMock.On("Send", actual, resp)

	actual.Send(resp)

	botMock.AssertExpectations(t)
	botMock.AssertNumberOfCalls(t, "Send", 1)
}
