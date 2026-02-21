package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
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

		content, err := generateArticleHtml(article, domain.ArticleTemplate.String())
		if err != nil {
			return fmt.Errorf("error generating html file content: %w", err)
		}
		article.Content = content

		articles[articleLanguage] = append(articles[articleLanguage], article)
		return nil
	})
	if err != nil {
		fmt.Println("error while walking through blog file path: %w", err)
		return
	}

	// Order in Most recent first
	for i := range domain.LanguageCount {
		slices.Reverse(articles[i])
	}

	fmt.Println(articles)
	// I need to add the spanish version though
	// Think of a way to create the content for the specified file. I think the best way is just to pass the file and that's it. We already have the file type in it
	index := domain.NewFile(domain.English, "index.html", domain.Index)
	index.Content = fmt.Sprintf(domain.HomePage.String(), generateIndexHtml(articles, projects))

	articlesIndex := domain.NewFile(domain.English, "index.html", domain.Article)
	articlesIndex.Content := fmt.Sprintf(domain.ArticlesPage.String(), domain.Navbar.String(), generateArticlesIndexHtml(articles))

	files := append(articles, index, articlesIndex)

	// Upload
	if !flags.IsSet(flags.Debug) {
		InitUpload(context.Background(), files)
	}
	
	// Debug
	// if flags.IsSet(flags.Debug) {
	// 	InitDebug(context.Background(), files, domain.PublicProjects)
	// }
}
