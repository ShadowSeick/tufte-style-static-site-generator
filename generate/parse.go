package generate

import (
	"strings"
)

type State uint8

const (
	Title State = iota
	NewSection
	InsideSection
	StateCount
)

var parseOrder = []MarkdownElement{Italic, Bold, InlineCode, Link, SideNote, MarginNote}

func parseLine(line string) (string, State) {
	if line == "" {
		return "", InsideSection
	}

	result, match := Header.Html(line)
	if match {
		headerNumber := strings.Count(string(line), "#")

		var state State
		switch headerNumber {
		case 1:
			state = Title
		case 2:
			state = NewSection
		case 3:
			state = InsideSection
		}

		return result, state
	}

	result, match = Subheader.Html(line)
	if match {
		return result, Title
	}

	result, match = Code.Html(line)
	if match {
		return result, InsideSection
	}

	result = line
	for _, markdownEl := range parseOrder {
		result, _ = markdownEl.Html(result)
	}

	result, match = Image.Html(result)
	if match {
		return result, InsideSection
	}

	result, _ = Paragraph.Html(result)

	return result, InsideSection
}
