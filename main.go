package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
	"github.com/ShadowSeick/tufte-style-static-site-generator/generate"
	"github.com/ShadowSeick/tufte-style-static-site-generator/pkg/flags"
)

func main() {
	// Get flags
	flags.Init()

	// Build Articles
	blogPath, err := filepath.Abs(domain.Article.LocalDirectory())
	if err != nil {
		fmt.Println("error getting working directory: %w", err)
		return
	}

	articles := make(map[domain.Language][]domain.File)
	var articleLanguage domain.Language
	err = filepath.WalkDir(blogPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		baseName := info.Name()
		if baseName == domain.Article.LocalDirectory() {
			return nil
		}

		if info.IsDir() && baseName == "assets" {
			return fs.SkipDir
		}

		if info.IsDir() {
			dirInfo, err := info.Info()
			if err != nil {
				return fmt.Errorf("error getting directory info: %w", err)
			}
			language, err := domain.ParseLanguage(dirInfo.Name())
			if err != nil {
				return fmt.Errorf("error parsing language: %w", err)
			}
			articleLanguage = language
			return nil
		}

		fileInfo, err := info.Info()
		if err != nil {
			return fmt.Errorf("error getting directory info: %w", err)
		}

		article, err := domain.NewFile(articleLanguage, fileInfo.Name(), domain.Article)
		if err != nil {
			return fmt.Errorf("error creating new article: %w", err)
		}

		content, err := generate.ArticleHtml(article)
		if err != nil {
			return fmt.Errorf("error generating html file content: %w", err)
		}
		article.Content = content

		articles[articleLanguage] = append(articles[articleLanguage], article)
		return nil
	})
	if err != nil {
		fmt.Println("error while walking through blog file path: ", err)
		return
	}

	// Order in Most recent first
	var files []domain.File
	for i := range domain.LanguageCount {
		slices.Reverse(articles[i])
		files = append(files, articles[i]...)
	}

	for i := range domain.LanguageCount {
		// Index page
		index, err := domain.NewFile(i, "index.html", domain.Index)
		if err != nil {
			fmt.Println("error creating new file: ", err)
			return
		}

		content, err := generate.IndexHtml(i, articles[i], domain.PublicProjects[i])
		if err != nil {
			fmt.Println("error generating index html: ", err)
			return
		}
		index.Content = content
		files = append(files, index)

		// Articles index page
		articlesIndex, err := domain.NewFile(i, "index.html", domain.ArticleIndex)
		if err != nil {
			fmt.Println("error creating new file: ", err)
			return
		}

		content, err = generate.ArticlesIndexHtml(i, articles[i])
		if err != nil {
			fmt.Println("error generating articles index html: ", err)
			return
		}
		articlesIndex.Content = content
		files = append(files, articlesIndex)

		// Contact page
		contactIndex, err := domain.NewFile(i, "contact.html", domain.ContactIndex)
		if err != nil {
			fmt.Println("error creating new file: ", err)
			return
		}

		contactIndex.Content = generate.ContactIndexHtml(i)
		files = append(files, contactIndex)
	}

	// Upload
	if !flags.IsSet(flags.Debug) {
		InitUpload(context.Background(), files)
	}

	// Debug
	if flags.IsSet(flags.Debug) {
		InitDebug(context.Background(), files)
	}
}
