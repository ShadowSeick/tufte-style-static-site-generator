package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidArticleBaseName = errors.New("invalid article base name")
)

const (
  ArticleFolder = "articles"
	ArticleDebugFolder = "debug"
	ArticleHtmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="/assets/styles/tufte.css"/>
		<script src="/assets/scripts/highlight/highlight.js"></script>
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


type Article struct {
	Name string
	Language Language
}

func (a *Article) FilePath() string {
	path, err := filepath.Abs(filepath.Join(ArticleFolder, a.Language.String(), a.Name))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (a *Article) StorageFilePath(directory string) string {
	return fmt.Sprintf("%s/%s/%s.%s", directory, a.Language, a.Name, HTML.String())
}

func (a *Article) Title() string {
	name := strings.Join(strings.Split(a.Name, "-"), " ")
	if len(name) == 0 {
		panic("article name should be at least 1 character long")
	}
	return strings.ToUpper(string(name[0])) + name[1:]
}

func NewArticle(language Language, fileBaseName string) (Article, error) {
	var article Article
	baseName := strings.Split(fileBaseName, ".")
	if len(baseName) != 2 {
		return article, fmt.Errorf("invalid article name: %s", fileBaseName)
	}

	article.Name = baseName[0]
	article.Language = language
	return article, nil
}
