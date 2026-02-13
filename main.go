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
	blogPath, err := filepath.Abs(ArticleFolders)
	if err != nil {
		fmt.Println("error getting working directory: %w", err)
		return
	}

	var articles []Article
	err = filepath.WalkDir(blogPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var article Article
		baseName := info.Name()
		if baseName == ArticleFolders {
			return nil
		}

		if info.IsDir() && baseName == "assets" {
			return fs.SkipDir
		}

		if info.IsDir() {
			article.name = info.Name()
			articles = append(articles, article)
			return nil
		}

		article = articles[len(articles) - 1]
		article.ParseBaseName(baseName)

		return nil
	})
	if err != nil {
		fmt.Println("error while walking through blog file path: %w", err)
		return
	}

	bunny.Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, article := range articles {
		content, err := GenerateHTML(article.FilePath(), ArticleHtmlTemplate)
		if err != nil {
			fmt.Println("error while generating html file content: ", err)
			return
		}
		checksum := strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(content))))

		files, err := bunny.GetFile(ctx, article.StorageFilePath())
		if err != nil {
			fmt.Println("error while getting file", err)
			return
		}

		var hasBeenUploaded bool
		for _, file := range files {
			if file.Checksum != nil && *file.Checksum == checksum {
				hasBeenUploaded = true
				break
			}
		}

		if !hasBeenUploaded {
			if err := bunny.UploadFile(ctx, article.StorageFilePath(), checksum, []byte(content)); err != nil {
				fmt.Println("error uploading file", err)
			}
		}
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
	
}
