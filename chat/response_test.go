package chat

import (
	"testing"

	assertion "github.com/stretchr/testify/assert"
)

func TestGetLanguageFromText(t *testing.T) {
	tt := []struct {
		title    string
		in       string
		language Language
		err      string
	}{
		{
			"defined language",
			"EN",
			EN,
			"",
		},
		{
			"undefined language",
			"UA",
			-1,
			"unknown language",
		},
	}

	for _, tc := range tt {
		t.Run(tc.title, func(t *testing.T) {
			assert := assertion.New(t)

			language, err := GetLanguageFromText(tc.in)

			assert.Equal(tc.language, language)

			if err != nil {
				assert.EqualError(err, tc.err)
			}
		})
	}

}

func TestCheckIfPartOfSpeechExists(t *testing.T) {
	tt := []struct {
		title        string
		language     Language
		partOfSpeech string
		isOk         bool
	}{
		{
			"defined language and part of speech",
			EN,
			"noun",
			true,
		},
		{
			"undefined language",
			Language(-1),
			"noun",
			false,
		},
		{
			"undefined part of speech",
			EN,
			"xxx",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.title, func(t *testing.T) {
			assert := assertion.New(t)

			isOk := CheckIfPartOfSpeechExists(tc.language, tc.partOfSpeech)
			assert.Equal(tc.isOk, isOk)
		})
	}
}
