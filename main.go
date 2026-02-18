package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/bunny"
	"github.com/ShadowSeick/tufte-style-static-site-generator/pkg/flags"
)

func main() {
	// Get flags
	flags.Init()
	
	// Build Articles
	blogPath, err := filepath.Abs(domain.ArticleFolder)
	if err != nil {
		fmt.Println("error getting working directory: %w", err)
		return
	}

	articles := make(map[domain.Language][]domain.Article)
	var articleLanguage domain.Language
	err = filepath.WalkDir(blogPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		baseName := info.Name()
		if baseName == domain.ArticleFolder {
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
		article, err := domain.NewArticle(articleLanguage, fileInfo.Name())
		if err != nil {
			return fmt.Errorf("error creating new article: %w", err)
		}

		articles[articleLanguage] = append(articles[articleLanguage], article)
		return nil
	})
	if err != nil {
		fmt.Println("error while walking through blog file path: %w", err)
		return
	}

	if !flags.IsSet(flags.Debug) {

		bunny.Init()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		uploadedFiles := make(map[domain.Language][]bunny.File)
		for i := range domain.LanguageCount {
			files, err := bunny.GetFile(ctx, fmt.Sprintf("%s/%s/", domain.ArticleFolder, i.String()))
			if err != nil {
				fmt.Println("error while getting file", err)
				return
			}

			uploadedFiles[i] = files
		}

		if err := uploadArticles(ctx, articles, uploadedFiles); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Finished uploading articles")
	}
	
	// Debug
	if flags.IsSet(flags.Debug) {
		if err := os.MkdirAll("./debug/en/", 0777); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating english articles directory")
		}

		if err := os.MkdirAll("./debug/es/", 0777); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating spanish articles directory")
		}

		if err := generateHTMLFile(articles, domain.ArticleDebugFolder); err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		listenToArticleChanges(articles, cancel)
		go regenerateHTMLFiles(ctx, articles)

		// Run server
		fs := http.FileServer(http.Dir(domain.ArticleDebugFolder))
		http.Handle("/", fs)

		log.Print("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}



func uploadArticles(ctx context.Context, articlesByLanguage map[domain.Language][]domain.Article, uploadedFiles map[domain.Language][]bunny.File) error {
	for _, articles := range articlesByLanguage {
		for _, article :=	range articles {
			content, err := generateHTML(article, domain.ArticleHtmlTemplate)
			if err != nil {
				return fmt.Errorf("error generating html file content: %w", err)
			}
			checksum := strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(content))))

			var hasBeenUploaded bool

			files := uploadedFiles[articel.Language]
			for _, file := range files {
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
	}
	return nil
}
