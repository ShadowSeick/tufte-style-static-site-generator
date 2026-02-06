package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

const (
	template = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="assets/tufte.css"/>
		<meta name="viewport" content="width=device-width, initial-scale=1">
	</head>
	<body>
		<article>
		%s
		</article>
	</body>
	</html>
	`
)

type MarkdownElement uint8

const (
	Header MarkdownElement = iota
	Subheader
	Paragraph
	Link
	SideNote
	MarginNote
	Code
	InlineCode
	Image
	Italic
	Bold
	MarkdownCount
)

var htmlString = [MarkdownCount]string{
	Header:     `<h%d id="%s">%s</h%d>`,
	Subheader:  `<p class="subtitle">%s</p>`,
	Paragraph:  `<p>%s</p>`,
	Link:       `<a href="%s">%s</a>`,
	SideNote:   `<label for="sn-%s" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-%s" class="margin-toggle"><span class="sidenote">%s</span>`,
	MarginNote: `<label for="mn-%s" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-%s" class="margin-toggle"><span class="marginnote">%s</span>`,
	Code:       `<pre><code>%s</code></pre>`,
	InlineCode: `<code>%s</code>`,
	Image:      `<figure>%s<img src="%s" alt="%s"/></figure>`,
	Italic:     `<em>%s</em>`,
	Bold:       `<b>%s</b>`,
}

var markdownRegexp = [MarkdownCount]*regexp.Regexp{
	Header:     regexp.MustCompile(`^#+\s+(.+)\s+\{@(.+)\}$`),
	Subheader:  regexp.MustCompile(`\[\^sub-header]\((.+?)\)`),
	Paragraph:  nil,
	Link:       regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`),
	SideNote:   regexp.MustCompile(`\[\^side-note\]\((.+?)\)`),
	MarginNote: regexp.MustCompile(`\[\^margin-note\]\((.+?)\)`),
	Code:       regexp.MustCompile("```([a-z]*)\\n([\\s\\S]+?)```"),
	InlineCode: regexp.MustCompile("`([^`]+)`"),
	Image:      regexp.MustCompile(`!\[([^\]]*)\]\(([^\)]+)\)`),
	Italic:     regexp.MustCompile(`\*\*(.+?)\*\*`),
	Bold:       regexp.MustCompile(`\*(.+?)\*`),
}

func (el MarkdownElement) Html(args ...any) string {
	if el >= MarkdownCount {
		panic("invalid markdown element")
	}
	return fmt.Sprintf(htmlString[el], args...)
}

func (el MarkdownElement) Match(text string) []string {
	if el >= MarkdownCount {
		panic("invalid markdown element")
	}
	return markdownRegexp[el].FindStringSubmatch(text)
}

func (el MarkdownElement) Text(matches []string) (string, string) {
	switch el {
	case Header:
		return matches[1], matches[2]
	case Subheader:
		return matches[1], ""
	case Link:
		return matches[1], matches[2]
	case SideNote:
		return matches[1], ""
	case MarginNote:
		return matches[1], ""
	case Code:
		return matches[1], matches[2]
	case InlineCode:
		return matches[1], ""
	case Image:
		return matches[1], matches[2]
	case Italic:
		return matches[1], ""
	case Bold:
		return matches[1], ""
	default:
		panic("element not handled")
	}
}

type State uint8

const (
	Title State = iota
	NewSection
	InsideSection
	StateCount
)

func ParseLine(line string) (string, State) {
	// Need only to check at first
	var state State
	if line == "" {
		return "", InsideSection
	}

	// HEADER
	matches := Header.Match(line)
	if matches != nil {
		headerNumber := strings.Count(string(line), "#")
		header, id := Header.Text(matches)

		switch headerNumber {
		case 1:
			state = Title
		case 2:
			state = NewSection
		case 3:
			state = InsideSection
		}

		return Header.Html(headerNumber, id, header, headerNumber), state
	}

	// SUBHEADER
	matches = Subheader.Match(line)
	if matches != nil {
		subheader, _ := Subheader.Text(matches)
		return Subheader.Html(subheader), Title
	}

	// IMAGE
	matches = Image.Match(line)
	if matches != nil {
		alt, path := Image.Text(matches)

		// Margin note
		var marginNote string
		marginMatch := MarginNote.Match(line)
		// Need to match for bold, italic, link and inline code
		if marginMatch != nil {
			// Check for bold and italic
			marginNote, _ = MarginNote.Text(marginMatch)
		}
		return Image.Html(MarginNote.Html(marginNote), path, alt), InsideSection
	}
	
	// What can go to a paragraph?
	// - Sidenote

	// What can go inside everything?
	// - Link
	// - Inline code
	// - Bold
	// - Italic

	// Code -- How it is done in Markdwon
	// Paragraph
	// Search for link, inline code, sidenote

	return fmt.Sprintf(htmlString[Paragraph], line), InsideSection
}

// Maybe this approach is too naive. The problem is that sometimes I like to add Introduction while other I don't. I need to put both and the only way is to build a structure around into
// I need to make some building blocks around them, but that would complicate a lot things... I still need to thing how to do it.

func main() {
	file, err := os.OpenFile("./example.md", os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Println("error opening the file: %w", err)
		return
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
				return
			}
			break
		}

		if startingCodeText := strings.Index(line, "```"); startingCodeText >= 0 {
			code, err := reader.ReadString('`')
			if err != nil {
				fmt.Println("error trying to read code:", err)
				return
			}

			line += code
			line += "```"
			reader.Discard(4) // Dicard last 3 bytes ``\n
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

	fmt.Print(html.String())
}
