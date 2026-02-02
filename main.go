package main

import (
    "bufio"
    "fmt"
    "os"
		"io/fs"
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

var htmlString = [MarkdownCount]string {
	Header: `<h%d id="%s">%s</h%d>`,
	Subheader: `<p class="subtitle">%s</p>`,
	Paragraph: `<p>%s</p>`,
	Link: `<a href="%s">%s</a>`,
	SideNote: `<label for="sn-%s" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-%s" class="margin-toggle"><span class="sidenote">%s</span>`,
	MarginNote: ``,
	Code: `<pre><code>%s</code></pre>`,
	InlineCode: `<code>%s</code>`,
	Image: `<figure src="%s">%s</figure>`,
	Italic: `<em>%s</em>`,
	Bold: `<b>%s</b>`,
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
	// HEADER
	var state State
	if line == "" {
		return "", InsideSection
	}
	headerNumber := strings.Count(string(line), "#")
	if headerNumber > 0 {
		// Get header
		header := line[headerNumber+1:]
		// Check for the id
		var id string
		startIndexID := strings.IndexRune(line, '{')
		if startIndexID > headerNumber {
			// Get the actual id
			id = header[startIndexID-1:len(header)-1]
			header = header[:startIndexID-3]
		}

		switch headerNumber {
		case 1:
			state = Title
		case 2:
			state = NewSection
		case 3:
			state = InsideSection
		}

		return fmt.Sprintf(htmlString[Header], headerNumber, id, header, headerNumber), state
	}

	// Subheader @subheader
	
	
	// What can go to a paragraph?
	// - Link
	// - Inline code
	// - Sidenote
	
	// What can go into images?
	// - Margin note

	// What will never go alone?
	// - Side notes
	// - Margin notes
	
	// Side note [*side-note](side note)
	// Margin note [^margin-note](margin note)
	// Code -- How it is done in Markdwon
	// Image -- ![alt-text](filepath)
	// Paragraph
	return fmt.Sprintf(htmlString[Paragraph], line), InsideSection
}


// Maybe this approach is too naive. The problem is that sometimes I like to add Introduction while other I don't. I need to put both and the only way is to build a structure around into
// I need to make some building blocks around them, but that would complicate a lot things... I still need to thing how to do it.

func main() {
	file, err := os.OpenFile("./example.md", os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Errorf("error opening the file: %w", err)
		return
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)

	oldState := Title
	var html strings.Builder
	for scanner.Scan() {
		// name := file.Name()
		// Generate the proper HTML page from markdown
		// Need to keep track of where am I. For that I can make some sort of state machine
		
		text, newState := ParseLine(scanner.Text())

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
