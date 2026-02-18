package main

import (
	"context"
	"path/filepath"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
)

type articleChanged struct {
	index int
	language domain.Language
}

var (
	changesChannel = make(chan articleChanged)
	initialize sync.Once
)

func InitDebug(ctx context.Context, articlesByLanguage map[domain.Language][]domain.Article) {
	initialize.Do(func() {
		cwd, _ := os.Getwd()
		if err := os.MkdirAll(filepath.Join(cwd, domain.ArticleDebugFolder, domain.English.String()), 0755); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating english articles directory %v", err)
		}

		if err := os.MkdirAll(filepath.Join(cwd, domain.ArticleDebugFolder, domain.Spanish.String()), 0755); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating spanish articles directory")
		}

		newCtx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		for _, articles := range articlesByLanguage {
			for index, article := range articles {
				if err := generateHTMLFile(article); err != nil {
					log.Fatal(err)
				}

				// Hot reload articles
				wg.Go(func () {
					defer wg.Done()
					listenToArticleChanges(newCtx, index, article, cancel)
				})

				wg.Go(func () {
					defer close(changesChannel)
					wg.Wait()
				})
			}
		}

		go regenerateHTMLFiles(ctx, articlesByLanguage)

		// Run server
		mux := http.NewServeMux()

		articles := http.FileServer(http.Dir(domain.ArticleDebugFolder))
		mux.Handle("/articles/", http.StripPrefix("/articles/", articles))

		assets := http.FileServer(http.Dir("assets"))
		mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

		
		log.Print("Listening on :3000...")
		log.Fatal(http.ListenAndServe(":3000", mux))
	})
}

func regenerateHTMLFiles(ctx context.Context, articlesByLanguage map[domain.Language][]domain.Article) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("gracefully shutting down")
			return
		case article, ok := <-changesChannel:
			if !ok { // Channel closed
				return
			}

			if err := generateHTMLFile(articlesByLanguage[article.language][article.index]); err != nil {
				fmt.Println("error regenerating html file", articlesByLanguage[article.language][article.index].Name, err)
				return
			}
		}
	}
}

func listenToArticleChanges(ctx context.Context, index int, article domain.Article, cancel context.CancelFunc) {
	filepath := article.FilePath()
	currStat, err := os.Stat(filepath)
	if err != nil {
		fmt.Println("error has occurred when listening to file stat", filepath, err)
		cancel()
		return
	}
	ticker := time.NewTicker(5*time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			fmt.Println("gracefully shuttinwdown listening to article: ", article.Name)
			return
		case <-ticker.C:
			stat, err := os.Stat(filepath)
			if err != nil {
				fmt.Println("error has occurred when listening to file stat", filepath, err)
				cancel()
				return
			}

			if stat.Size() != currStat.Size() || !stat.ModTime().Equal(currStat.ModTime()) {
				currStat = stat
				changesChannel <- articleChanged{
					index: index,
					language: article.Language,
				}
			}
		}
	}
}

func generateHTMLFile(article domain.Article) error {
	content, err := generateHTML(article, domain.ArticleHtmlTemplate)
	if err != nil {
		return fmt.Errorf("error generating html file content: %w", err)
	}

	if err := os.WriteFile(article.StorageFilePath(domain.ArticleDebugFolder), []byte(content), 0644); err != nil {
		return fmt.Errorf("error creating html file: %w", err)
	}
	return nil
}
