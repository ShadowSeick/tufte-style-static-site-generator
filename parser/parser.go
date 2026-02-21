package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"bufio"
	"io/fs"
	"fmt"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
)

type State uint8

const (
	Title State = iota
	NewSection
	InsideSection
	StateCount
)

var parseOrder = []domain.MarkdownElement{domain.Italic, domain.Bold, domain.InlineCode, domain.Link, domain.SideNote, domain.MarginNote}
func ParseLine(line string) (string, State) {
	if line == "" {
		return "", InsideSection
	}

	result, match := domain.Header.Html(line)
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

	result, match = domain.Subheader.Html(line)
	if match {
		return result, Title
	}

	result, match = domain.Code.Html(line)
	if match {
		return result, InsideSection
	}

	result = line
	for _, markdownEl := range parseOrder {
		result, _ = markdownEl.Html(result)
	}

	result, match = domain.Image.Html(result)
	if match {
		return result, InsideSection
	}

	result, _ = domain.Paragraph.Html(result)

	return result, InsideSection
}

// I would need only to expose a GenerateHtml function that can actually choose how to parse it
// Actually think of a better way to deal with it, I need to pass other things depending on the type they are, so I think I should have 3 functions
// - Article
// - Index
// - ArticleIndex
// Probably not needed and I can do it with just one
// Let's continue in the next session
func GenerateHtml(file domain.File) (string, error) {
	if file.Type != domain.Article || file.Type != domain.Index || file.Type != domain.ArticlesIndex {
		return "", fmt.Errorf("wrong file type")
	}

	file, err := os.OpenFile(file.LocalFilePath(), os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Println("error opening the file: %w", err)
		return "", err
	}
	defer file.Close()

	switch file.Type {
		// Here I need to watch out, how do I deal with ArticleIndex?
	case domain.Article:
		return generateArticleHtml(file)
	case domain.ArticlesIndex:
		return generateArticlesIndexHtml(file)
	case domain.Index:
		return generateIndexHtml(file)
	default:
		return "", fmt.Errorf("not handled file type")
	}
}

func GenerateArticlesIndexHtml()

func generateArticleHtml(file *os.File) (string, error) {
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

	return fmt.Sprintf(ArticlePage.String(), article.Title(), domain.Navbar.String(), html.String()), nil
}

func generateIndexHtml(articles []domain.File, projects []domain.Project) string {
	var html strings.Builder
	html.WriteString(domain.Navbar.String())
	html.WriteString(domain.ContactInfo.String())

	var projectsInfo string
	for _, project := range projects {
		projectRow := fmt.Sprintf(domain.Td.String(), fmt.Sprintf(domain.A.String(), project.Link, project.Name))
		projectRow += fmt.Sprintf(domain.Td.String(), fmt.Sprintf(domain.P.String(), project.Description))
		projectsInfo += fmt.Sprintf(domain.Tr.String(), projectRow)
	}
	html.WriteString(fmt.Sprintf(domain.ProjectsInfo.String(), projectsInfo))

	var articlesList string
	for _, article := range articles {
		articleItem := fmt.Sprintf(domain.H3.String(), fmt.Sprintf(domain.A.String(), "/"+article.StorageFilePath(domain.ArticleFolder), article.Title()))
		articleItem += fmt.Sprintf(domain.P.String(), article.DateString())
		articlesList += fmt.Sprintf(domain.Li.String(), articleItem)
	}
	html.WriteString(fmt.Sprintf(domain.ArticlesInfo.String(), articlesList))
	return html.String()
}

func generateArticlesIndexHtml(articles []domain.Article) string {
	var html strings.Builder

	var articlesList string
	for _, article := range articles {
		articleItem := fmt.Sprintf(domain.H3.String(), fmt.Sprintf(domain.A.String(), "/"+article.StorageFilePath(domain.ArticleFolder), article.Title()))
		articleItem += fmt.Sprintf(domain.P.String(), article.DateString())
		articlesList += fmt.Sprintf(domain.Li.String(), articleItem)
	}

	html.WriteString(articlesList)
	return html.String()
}
