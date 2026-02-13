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
  ArticleFolders = "blog"
	ArticleHtmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="assets/styles/tufte.css"/>
		<script src="assets/scripts/highlight.js"></script>
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
	checksum string
}

func (a *Article) FilePath() string {
	path, err := filepath.Abs(filepath.Join(ArticleFolders, a.name, a.language.String()))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (a *Article) StorageFilePath() string {
	return fmt.Sprintf("articles/%s/%s.%s", a.language, a.name, HTML.String())
}

// func (a *Article) HtmlFilePath(tempDirPath string) string {
// 	return fmt.Sprintf("%s-%s.%s", filepath.Join(tempDirPath, a.name), a.language, HTML.String())
// }

func (a *Article) ParseBaseName(baseName string) error {
	articleFileInfo := strings.Split(baseName, ".")
	if len(articleFileInfo) != 2 {
		return ErrInvalidArticleBaseName
	}

	language, err := ParseLanguage(articleFileInfo[0])
	if err != nil {
		return err
	}
	a.language = language

	return nil
}
