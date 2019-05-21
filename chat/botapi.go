package chat

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BotAPI interface {
	Send(chat *Chat, response *Responses)
}

type Bot struct {
	*tgbotapi.BotAPI
}

func (bot *Bot) Send(chat *Chat, response *Responses) {
	msg := tgbotapi.NewMessage(chat.State.ChatID, response.Text)

	if response.ReplyKeyboardMarkup != nil {
		msg.ReplyMarkup = response.ReplyKeyboardMarkup
	}

	if response.ReplyKeyboardMarkup == nil {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	if _, err := bot.BotAPI.Send(msg); err != nil {
		chat.Logger.Error(err)
	}
}
