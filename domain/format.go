package domain

import (
	"errors"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

type Format uint8

const (
	Markdown Format = iota
	HTML
	FormatCount
)

var formatNames = [FormatCount]string{
	Markdown: "md",
	HTML: "html",
}

func (f Format) String() string {
	if f >= FormatCount {
		panic(ErrInvalidFormat.Error())
	}
	return formatNames[f]
}
