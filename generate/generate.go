package generate

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

// I actually think this could be unified to just 1 function, but it would overcomplicate things
func ArticleHtml(file domain.File) (string, error) {
	if file.Type != domain.Article {
		return "", fmt.Errorf("wrong file type")
	}

	osFile, err := os.OpenFile(file.LocalFilePath(), os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		fmt.Println("error opening the file: %w", err)
		return "", err
	}
	defer osFile.Close()

	reader := bufio.NewReader(osFile)

	oldState := Title
	var html strings.Builder

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

		text, newState := parseLine(line[:len(line)-1])

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

	return fmt.Sprintf(ArticleTemplate.String(file.Language), file.Title(), Navbar.String(file.Language), html.String()), nil
}

// This is comboluted, I don't like it. This would be much simpler with a strings builder from scratch
func IndexHtml(language domain.Language, files []domain.File, projects []domain.Project) (string, error) {
	for _, file := range files {
		if file.Type != domain.Article {
			return "", fmt.Errorf("wrong file type: index only accepts article files")
		}
	}

	var html strings.Builder
	html.WriteString(Navbar.String(language))

	var projectsInfo string
	for _, project := range projects {
		projectRow := fmt.Sprintf(Td.String(), fmt.Sprintf(A.String(), project.Link, project.Name))
		projectRow += fmt.Sprintf(Td.String(), fmt.Sprintf(P.String(), project.Description))
		projectsInfo += fmt.Sprintf(Tr.String(), projectRow)
	}
	html.WriteString(fmt.Sprintf(ProjectsInfo.String(language), projectsInfo))

	var articlesList string
	for _, file := range files {
		articleItem := fmt.Sprintf(H3.String(), fmt.Sprintf(A.String(), "/"+file.RemoteFilePath(), file.Title()))
		articleItem += fmt.Sprintf(P.String(), file.DateString())
		articlesList += fmt.Sprintf(Li.String(), articleItem)
	}
	html.WriteString(fmt.Sprintf(ArticlesInfo.String(language), articlesList))
	return fmt.Sprintf(HomePage.String(language), html.String()), nil
}

func ArticlesIndexHtml(language domain.Language, files []domain.File) (string, error) {
	for _, file := range files {
		if file.Type != domain.Article {
			return "", fmt.Errorf("wrong file type: articles index only accepts article files")
		}
	}

	var html strings.Builder
	var articlesList string
	for _, file := range files {
		articleItem := fmt.Sprintf(H3.String(), fmt.Sprintf(A.String(), "/"+file.RemoteFilePath(), file.Title()))
		articleItem += fmt.Sprintf(P.String(), file.DateString())
		articlesList += fmt.Sprintf(Li.String(), articleItem)
	}

	html.WriteString(articlesList)
	return fmt.Sprintf(ArticlesPage.String(language), Navbar.String(language), html.String()), nil
}

func ContactIndexHtml(language domain.Language) string {
	return fmt.Sprintf(ContactPage.String(language), Navbar.String(language))
}
