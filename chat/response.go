package chat

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Responses struct {
	Text                 string
	InlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
	ReplyKeyboardMarkup  *tgbotapi.ReplyKeyboardMarkup
	ParseMode            string
}

type Language int

const (
	EN Language = iota
	PL
)

func GetLanguageFromText(text string) Language {
	switch text {
	case "EN":
		return EN
	case "PL":
		return PL
	default:
		return EN
	}
}

type Message string

const (
	Welcome                Message = "Hey! Seems you have new words.\nLet's save it!"
	UnknownCommand         Message = "Unknown command."
	OnReceivedWord         Message = "Got you word. Pick language to save this word."
	OnReceivedLanguage     Message = "Nice, then choose part of speech."
	OnReceivedPartOfSpeech Message = "Type sentences to see how to use your word"
	OnReceivedSentences    Message = "Brilliant! All info has been saved.\nIf you need to save another word, just type it."
)

func LanguageKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("EN"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("PL"),
		),
	)

	return &keyboard
}

func PartOfSpeech(language Language) *tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	if language == EN {
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("noun"),
				tgbotapi.NewKeyboardButton("verb"),
				tgbotapi.NewKeyboardButton("adjective"),
				tgbotapi.NewKeyboardButton("preposition"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("adverb"),
				tgbotapi.NewKeyboardButton("pronoun"),
				tgbotapi.NewKeyboardButton("conjunction"),
				tgbotapi.NewKeyboardButton("interjection"),
			),
		)
	}

	if language == PL {
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("rzeczownik"),
				tgbotapi.NewKeyboardButton("czasownik"),
				tgbotapi.NewKeyboardButton("przymiotnik"),
				tgbotapi.NewKeyboardButton("liczebnik"),
				tgbotapi.NewKeyboardButton("zaimek"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("przysłówek"),
				tgbotapi.NewKeyboardButton("przyimek"),
				tgbotapi.NewKeyboardButton("spójnik"),
				tgbotapi.NewKeyboardButton("wykrzyknik"),
				tgbotapi.NewKeyboardButton("partykuła"),
			),
		)
	}

	return &keyboard
}

func Send(bot *tgbotapi.BotAPI, response *Responses, chatState *State) {

	msg := tgbotapi.NewMessage(chatState.ChatID, response.Text)

	if response.InlineKeyboardMarkup != nil {
		msg.ReplyMarkup = response.InlineKeyboardMarkup
	}
	if response.ReplyKeyboardMarkup != nil {
		msg.ReplyMarkup = response.ReplyKeyboardMarkup
	}

	if response.ParseMode != "" {
		msg.ParseMode = response.ParseMode
	}

	if response.ReplyKeyboardMarkup == nil && response.InlineKeyboardMarkup == nil {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	if _, err := bot.Send(msg); err != nil {
		logger.Error(err)
	}
}