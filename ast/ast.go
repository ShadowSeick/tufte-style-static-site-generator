package ast

import (
	"fmt"
	"strings"
)

type astType uint8

const (
	header astType = iota
	subheader
	note
	link
	image
	code
	fontModifier
	paragraph
	astTypeCount
)

var astString = [astTypeCount]string{
	header:       "header",
	subheader:    "subheader",
	note:         "note",
	link:         "link",
	image:        "image",
	code:         "code",
	fontModifier: "font_modifier",
	paragraph:    "paragraph",
}

func (a astType) String() string {
	return astString[a]
}

type AST interface {
	HTML() string
	String() string
}

// HEADER

type Header struct {
	Level   int
	Content []byte
}

func (h Header) String() string {
	return fmt.Sprintf("%s: level %d; content: %s", header.String(), h.Level, h.Content)
}

func (h Header) HTML() string {
	return fmt.Sprintf("<h%[1]d>%s</h%[1]d", h.Level, h.Content)
}

// SUBHEADER

type Subheader struct {
	Content []byte
}

func (s Subheader) String() string {
	return fmt.Sprintf("%s: %s", subheader.String(), s.Content)
}

func (s Subheader) HTML() string {
	return fmt.Sprintf(`<p class="subtitle">%s</p>`, s.Content)
}

// NOTES

type noteType uint8

const (
	Margin noteType = iota
	Side
	noteTypeCount
)

var noteTypeString = [noteTypeCount]string{
	Margin: "margin",
	Side:   "side",
}

func (nt noteType) String() string {
	return noteTypeString[nt]
}

type Note struct {
	NoteType  noteType
	ClassName []byte
	Content   []byte
	Children  []AST
}

func (n Note) String() string {
	return fmt.Sprintf("%s %s: className %s", n.NoteType.String(), note.String(), n.ClassName)
}

func (n Note) HTML() string {
	var content strings.Builder
	for _, c := range n.Children {
		content.WriteString(c.HTML())
	}
	return content.String()
}
