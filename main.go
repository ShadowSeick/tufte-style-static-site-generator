package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/bunny"
)

func main() {
	// Build Articles
	blogPath, err := filepath.Abs(ArticleFolder)
	if err != nil {
		fmt.Println("error getting working directory: %w", err)
		return
	}

	var englishArticles []Article
	var spanishArticles []Article
	var articleLanguage Language
	err = filepath.WalkDir(blogPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		baseName := info.Name()
		if baseName == ArticleFolder {
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
			language, err := ParseLanguage(dirInfo.Name())
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
		article, err := NewArticle(articleLanguage, fileInfo.Name())
		if err != nil {
			return fmt.Errorf("error creating new article: %w", err)
		}

		switch articleLanguage {
		case English:
			englishArticles = append(englishArticles, article)
		case Spanish:
			spanishArticles = append(spanishArticles, article)
		}

		return nil
	})
	if err != nil {
		fmt.Println("error while walking through blog file path: %w", err)
		return
	}

	// Upload files
	bunny.Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var englishFiles []bunny.File
	var spanishFiles []bunny.File
	for i := range LanguageCount {
		files, err := bunny.GetFile(ctx, fmt.Sprintf("%s/%s/", ArticleFolder, i.String()))
		if err != nil {
			fmt.Println("error while getting file", err)
			return
		}

		switch i {
		case English:
			englishFiles =	files
		case Spanish:
			spanishFiles = files
		}
	}

	if err := UploadArticles(ctx, englishArticles, englishFiles); err != nil {
		fmt.Println(err)
		return
	}

	if err := UploadArticles(ctx, spanishArticles, spanishFiles); err != nil {
		fmt.Println(err)
		return
	}

	// Maybe it's a good idea for debugging. Do not delete it now
	// tempDir, err := os.MkdirTemp("", "temp")
	// if err != nil {
	// 	fmt.Println("error while creating temp dir for html files: ", err)
	// 	return
	// }
	//
	// for _, article := range articles {
	// 	content, err := GenerateHTML(article.FilePath(), ArticleHtmlTemplate)
	// 	if err != nil {
	// 		fmt.Println("error while generating html file content: ", err)
	// 		return
	// 	}
	//
	// 	path := article.HtmlFilePath(tempDir)
	// 	fmt.Println(path)
	// 	if err := os.WriteFile(article.HtmlFilePath(tempDir), []byte(content), 0644); err != nil {
	// 		fmt.Println("error writting html file: ", err)
	// 		return
	// 	}
	//
	// 	article.checksum = sha256.Sum256([]byte(content))
	// }
	
	fmt.Println("Finished uploading articles")
}

func UploadArticles(ctx context.Context, articles []Article, uploadedFiles []bunny.File) error {
	for _, article := range articles {
		content, err := GenerateHTML(article, ArticleHtmlTemplate)
		if err != nil {
			return fmt.Errorf("error generating html file content: %w", err)
		}
		checksum := strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(content))))

		var hasBeenUploaded bool
		for _, file := range uploadedFiles {
			if file.Checksum != nil && *file.Checksum == checksum {
				hasBeenUploaded = true
				break
			}
		}

		if !hasBeenUploaded {
			fmt.Println("uploading file...", article.Title())
			if err := bunny.UploadFile(ctx, article.StorageFilePath(), checksum, []byte(content)); err != nil {
				fmt.Println("error uploading file", err)
				return fmt.Errorf("error updating file: %w", err)
			}
		}
	}
	return nil
}
