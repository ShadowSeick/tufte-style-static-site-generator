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
		<link rel="stylesheet" href="assets/default.min.css">
		<script src="assets/highlight.min.js"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<script>hljs.highlightAll();</script>
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
	Subheader:  `<p class="subtitle">${text}</p>`,
	Paragraph:  `<p>%s</p>`,
	Link:       `<a href="${url}">${text}</a>`,
	SideNote:   `<label for="sn-%s" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-%s" class="margin-toggle"><span class="sidenote">%s</span>`,
	MarginNote: `<label for="mn-%s" class="margin-toggle">&#8853;</label><input type="checkbox" id="mn-%s" class="margin-toggle"><span class="marginnote">%s</span>`,
	Code:       `<pre><code class="language-${lang}">${code}</code></pre>`,
	InlineCode: `<code>${code}</code>`,
	Image:      `<img src="${url}" alt="${alt}"/>`,
	Italic:     `<em>${text}</em>`,
	Bold:       `<b>${text}</b>`,
}

var markdownRegexp = [MarkdownCount]*regexp.Regexp{
	Header:     regexp.MustCompile(`^(?P<hashes>#+)\s+(?P<text>.+)\s+\{@(?P<id>.+)\}$`),
	Subheader:  regexp.MustCompile(`\[\^sub-header\]\((?P<text>.+?)\)`),
	Paragraph:  nil,
	Link:       regexp.MustCompile(`\[(?P<text>[^\]]+)\]\((?P<url>[^\)]+)\)`),
	SideNote:   regexp.MustCompile(`\[\^side-note @(?P<id>[^\]]+)\]`),
	MarginNote: regexp.MustCompile(`\[\^margin-note @(?P<id>[^\]]+)\]`),
	Code:       regexp.MustCompile("```(?P<lang>[a-z]*)\\n(?P<code>[\\s\\S]+?)```"),
	InlineCode: regexp.MustCompile("`(?P<code>[^`]+)`"),
	Image:      regexp.MustCompile(`!\[(?P<alt>[^\]]*)\]\((?P<url>[^\)]+)\)`),
	Italic:     regexp.MustCompile(`\*\*(?P<text>.+?)\*\*`),
	Bold:       regexp.MustCompile(`\*(?P<text>.+?)\*`),
}

func (el MarkdownElement) Html(text string) (string, bool) {
	if el >= MarkdownCount {
		panic("invalid markdown element")
	}

	if el != Paragraph && !markdownRegexp[el].MatchString(text) {
		return text, false
	}

	result := text
	switch el {
	case Header:
		matches := markdownRegexp[Header].FindStringSubmatch(text)
		level := len(matches[1])
		id := matches[3]
		text := matches[2]

		result = fmt.Sprintf(`<h%d id="%s">%s</h%d>`, level, id, text, level)
	case SideNote, MarginNote: // This is fucking ugly but it works. Maybe there is a better way to do it later
		var newText strings.Builder

		endIndex := len(text) -1
		for _, match := range markdownRegexp[el].FindAllStringIndex(text, -1) {
			// Write until match
			newText.WriteString(text[:match[0]])

			// Get index if any
			startIndexID := strings.IndexRune(text[match[0]:match[1]], '@')
			var id string
			if startIndexID >= 0 {
				id = text[startIndexID+match[0]+1:match[1]-1]
				fmt.Println("ID:", id)
			}
			
			// Get Balanced parenthesis
			startParenthesis := match[1] + 1
			endParenthesis := len(text) - 1
			numberParenthesis := 1
			for pos, char := range text[startParenthesis:] {
				if char == '(' {
					numberParenthesis += 1
				}
				if char == ')' {
					numberParenthesis -= 1
				}
				if numberParenthesis == 0 {
					endParenthesis = startParenthesis + pos
					break
				}
			}

			newText.WriteString(fmt.Sprintf(htmlString[el], id, id, text[startParenthesis:endParenthesis]))
			endIndex = endParenthesis
		}

		if endIndex != len(text) -1 {
			newText.WriteString(text[endIndex+1:])
		}
		
		result = newText.String()
	case Image:
		result = fmt.Sprintf("<figure>%s</figure>", markdownRegexp[Image].ReplaceAllString(text, htmlString[Image]))
	case Paragraph:
		result = fmt.Sprintf(htmlString[el], result)
	default:
    result = markdownRegexp[el].ReplaceAllString(text, htmlString[el])
	}

	return result, true
}

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

	for _, markdownEl := range parseOrder {
		result, _ = markdownEl.Html(line)
	}

	result, match = Image.Html(line)
	if match {
		return result, InsideSection
	}

	result, _ = Paragraph.Html(result)

	return result, InsideSection
}

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

	fmt.Print(fmt.Sprintf(template, file.Name(), html.String()))
}
