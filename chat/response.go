package chat

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Responses struct {
	Text                string
	ReplyKeyboardMarkup *tgbotapi.ReplyKeyboardMarkup
}

type Language int

const (
	EN Language = iota
	PL
)

func GetLanguageFromText(text string) (Language, error) {
	switch text {
	case "EN":
		return EN, nil
	case "PL":
		return PL, nil
	default:
		return -1, errors.New("unknown language")
	}
}

type Message string

const (
	Welcome                Message = "Hey! Seems you have new words.\nLet's save it!"
	UnknownCommand         Message = "Unknown command."
	UnknownLanguage        Message = "Your language is not supported. Please, select from keyboard."
	UnknownPartOfSpeech    Message = "Please, select part of speech."
	OnReceivedWord         Message = "Got you word. Pick language to save this word."
	OnReceivedLanguage     Message = "Nice, then choose part of speech."
	OnReceivedPartOfSpeech Message = "Type sentences to see how to use your word."
	OnReceivedSentences    Message = "Brilliant! All info has been saved.\nIf you need to save another word, just type it."
)

var ENPartOfSpeech = []string{
	"noun",
	"verb",
	"adjective",
	"preposition",
	"adverb",
	"pronoun",
	"conjunction",
	"interjection",
}

var PLPartOfSpeech = []string{
	"rzeczownik",
	"czasownik",
	"przymiotnik",
	"liczebnik",
	"zaimek",
	"przysłówek",
	"przyimek",
	"spójnik",
	"wykrzyknik",
	"partykuła",
}

func CheckIfPartOfSpeechExists(language Language, partOfSpeech string) bool {
	var list []string

	if language == EN {
		list = ENPartOfSpeech
	}

	if language == PL {
		list = PLPartOfSpeech
	}

	for _, b := range list {
		if b == partOfSpeech {
			return true
		}
	}
	return false
}

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

func createPartOfSpeechKeyboard(inRow int, partOfSpeech []string) [][]tgbotapi.KeyboardButton {
	var keyboard [][]tgbotapi.KeyboardButton
	var buttons []tgbotapi.KeyboardButton

	for _, value := range partOfSpeech {
		button := tgbotapi.NewKeyboardButton(value)
		buttons = append(buttons, button)
	}

	for i := 0; i < len(partOfSpeech); i += inRow {
		end := i + inRow

		if end > len(partOfSpeech) {
			end = len(partOfSpeech)
		}

		keyboard = append(keyboard, buttons[i:end])
	}

	return keyboard
}

func PartOfSpeech(language Language) *tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	if language == EN {
		buttons := createPartOfSpeechKeyboard(4, ENPartOfSpeech)
		keyboard = tgbotapi.NewReplyKeyboard(
			buttons...,
		)
	}

	if language == PL {
		buttons := createPartOfSpeechKeyboard(4, PLPartOfSpeech)
		keyboard = tgbotapi.NewReplyKeyboard(
			buttons...,
		)
	}

	return &keyboard
}

func Send(context *Context, response *Responses) {
	msg := tgbotapi.NewMessage(context.State.ChatID, response.Text)

	if response.ReplyKeyboardMarkup != nil {
		msg.ReplyMarkup = response.ReplyKeyboardMarkup
	}

	if response.ReplyKeyboardMarkup == nil {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	if _, err := context.Bot.Send(msg); err != nil {
		context.Logger.Error(err)
	}
}
