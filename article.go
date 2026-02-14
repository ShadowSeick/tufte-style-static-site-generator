package main

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
	ArticleHtmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="/assets/styles/tufte.css"/>
		<script src="/assets/scripts/highlight.js"></script>
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
	name string
	language Language
}

func (a *Article) FilePath() string {
	path, err := filepath.Abs(filepath.Join(ArticleFolder, a.language.String(), a.name))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (a *Article) StorageFilePath() string {
	return fmt.Sprintf("%s/%s/%s.%s", ArticleFolder, a.language, a.name, HTML.String())
}

func (a *Article) Title() string {
	name := strings.Join(strings.Split(a.name, "-"), " ")
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

	article.name = baseName[0]
	article.language = language
	return article, nil
}

func (a *Article) HtmlFilePath(tempDirPath string) string {
	return fmt.Sprintf("%s-%s.%s", filepath.Join(tempDirPath, a.name), a.language, HTML.String())
}
