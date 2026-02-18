package domain

import (
	"errors"
)

var (
	ErrInvalidLanguage = errors.New("invalid language")
)

type Language uint8

const (
	English Language = iota
	Spanish
	LanguageCount
)

var languageNames = [LanguageCount]string{
	English: "en",
	Spanish: "es",
}

func (l Language) String() string {
	if l >= LanguageCount {
		panic(ErrInvalidLanguage.Error())
	}
	return languageNames[l]
}

func ParseLanguage(language string) (Language, error) {
	for i, l := range languageNames {
		if language == l {
			return Language(i), nil
		}
	}
	return LanguageCount, ErrInvalidLanguage
}
