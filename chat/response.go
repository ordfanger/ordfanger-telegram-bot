package chat

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Language type. Type alias for supported languages.
type Language int

// Message type. Used for defining bot's responses.
type Message string

// Supported languages.
const (
	EN Language = iota
)

// Bot's response messages.
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

// Responses used to define response text and keyboard for telegram bot.
type Responses struct {
	Text                string
	ReplyKeyboardMarkup *tgbotapi.ReplyKeyboardMarkup
}

// ENPartOfSpeech - defined part of speech for EN language.
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

// GetLanguageFromText returns language type from string.
func GetLanguageFromText(text string) (Language, error) {
	switch text {
	case "EN":
		return EN, nil
	default:
		return -1, errors.New("unknown language")
	}
}

// CheckIfPartOfSpeechExists checks if part of speech defined.
func CheckIfPartOfSpeechExists(language Language, partOfSpeech string) bool {
	var list []string

	if language == EN {
		list = ENPartOfSpeech
	}

	for _, b := range list {
		if b == partOfSpeech {
			return true
		}
	}
	return false
}

// LanguageKeyboard returns language keyboard.
func LanguageKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("EN"),
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

// PartOfSpeech returns part of speech keyboard.
func PartOfSpeech(language Language) *tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	if language == EN {
		buttons := createPartOfSpeechKeyboard(4, ENPartOfSpeech)
		keyboard = tgbotapi.NewReplyKeyboard(
			buttons...,
		)
	}

	return &keyboard
}
