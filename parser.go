package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"bufio"
	"io/fs"
	"fmt"
)

type State uint8

const (
	Title State = iota
	NewSection
	InsideSection
	StateCount
)

var parseOrder = []MarkdownElement{Italic, Bold, InlineCode, Link, SideNote, MarginNote}
func ParseLine(line string) (string, State) {
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

// Maybe this should be done directly passing the article
func GenerateHTML(article Article, template string) (string, error) {
	file, err := os.OpenFile(article.FilePath(), os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Println("error opening the file: %w", err)
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	oldState := Title
	var html strings.Builder
	// Instead of using this, use the raw read and break it into lines, this way I can proccess the exact same structure when I have some code
	var endOfLine bool
	for !endOfLine {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(io.EOF, err) {
				fmt.Println(fmt.Sprintf("error while reading string: %v", err))
				return "", err
			}
			break
		}

		if startingCodeText := strings.Index(line, "```"); startingCodeText >= 0 {
			code, err := reader.ReadString('`')
			if err != nil {
				fmt.Println("error trying to read code:", err)
				return "", err
			}

			line += code
			line += "```"
			reader.Discard(3) // Dicard last 3 bytes ```
		}

		text, newState := ParseLine(line[:len(line)-1])

		switch newState {
		case Title:
			html.WriteString(text)
		case NewSection:
			if oldState == InsideSection {
				html.WriteString("</section>")
			}
			html.WriteString("<section>")
			html.WriteString(text)
		case InsideSection:
			if oldState == Title {
				html.WriteString("<section>")
			}
			html.WriteString(text)
		default:
			panic("invalid state")
		}
		html.WriteString("\n")
		oldState = newState
	}
	html.WriteString("</section>")

	return fmt.Sprintf(template, article.Title(), html.String()), nil
}
