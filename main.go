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
		<meta name="viewport" content="width=device-width, initial-scale=1"
	</head>
	<body>
		<article>
		%s
		</article>
	</body>
	</html>
	`
	section = `<section>%s</section>`
)


type MarkdownElement uint8

const (
	Header MarkdownElement = iota
	Subheader
	Paragraph
	Link
	Sidenote
	Code
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
	Sidenote: `<label for="sn-%s" class="margin-toggle sidenote-number"></label><input type="checkbox" id="sn-%s" class="margin-toggle"><span class="sidenote">%s</span>`,
	Code: `<pre><code>%s</code></pre>`,
	Image: `<figure src="%s">%s</figure>`,
	Italic: `<em>%s</em>`,
	Bold: `<b>%s</b>`,
}

func ParseLine(l string) string {
	// Need only to check at first
	// HEADER
	headerNumber := strings.Count(string(l), "#")
	if headerNumber > 0 {
		// Get header
		header := l[headerNumber+1:]
		// Check for the id
		var id string
		startIndexID := strings.IndexRune(l, '{')
		if startIndexID > headerNumber {
			// Get the actual id
			id = header[startIndexID-1:len(header)-1]
			header = header[:startIndexID-3]
		}
		return fmt.Sprintf(htmlString[Header], headerNumber, id, header, headerNumber)
	}
	return ""

	// Subheader
	// Paragraph
	// Link
	// Sidenote
	// Code
	// Image
}



func main() {
	file, err := os.OpenFile("./example.md", os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Errorf("error opening the file: %w", err)
		return
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// name := file.Name()
		// Generate the proper HTML page from markdown
		s := ParseLine(scanner.Text())
		fmt.Println(s)
	}
}
