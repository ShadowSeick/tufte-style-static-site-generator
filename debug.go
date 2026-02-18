package main

import (
	"context"
	"fmt"
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
		if err := os.MkdirAll("./debug/en/", 0777); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating english articles directory")
		}

		if err := os.MkdirAll("./debug/es/", 0777); err != nil && !os.IsExist(err) {
			log.Fatalf("error creating spanish articles directory")
		}

		for 
		if err := generateHTMLFile(articles, domain.ArticleDebugFolder); err != nil {
			log.Fatal(err)
		}

		_, cancel := context.WithCancel(ctx)
		listenToArticleChanges(articles, cancel)
		go regenerateHTMLFiles(ctx, articles)

		// Run server
		fs := http.FileServer(http.Dir(domain.ArticleDebugFolder))
		http.Handle("/", fs)

		log.Print("Listening on :3000...")
		log.Fatal(http.ListenAndServe(":3000", nil))
	})
}

func regenerateHTMLFiles(ctx context.Context, articlesByLanguage map[domain.Language][]domain.Article) {
	for {
		select {
		case <-ctx.Done():
			return
		case article, ok := <-changesChannel:
			if !ok { // Channel closed
				return
			}

			if err := generateHTMLFile(articlesByLanguage[article.language][article.index], domain.ArticleDebugFolder); err != nil {
				fmt.Println("error regenerating html file", articlesByLanguage[article.language][article.index].Name, err)
				return
			}
		}
	}
}

func listenToArticleChanges(language domain.Language, articles []domain.Article, cancel context.CancelFunc) {
	defer close(changesChannel)

	var wg sync.WaitGroup
	for i, article := range articles {
		wg.Go(func () {
			defer wg.Done()
			filePath := article.FilePath()
			currStat, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("error has occurred when listening to file stat", filePath, err)
				cancel()
				return
			}
			
			for {
				stat, err := os.Stat(filePath)
				if err != nil {
					fmt.Println("error has occurred when listening to file stat", filePath, err)
					cancel()
					return
				}

				if stat.Size() != currStat.Size() || !stat.ModTime().Equal(currStat.ModTime()) {
					currStat = stat
					changesChannel <- articleChanged{
						index: i,
						language: language,
					}
				}

				time.Sleep(5 * time.Second)
			}
		})
	}
	wg.Wait()
}

func generateHTMLFile(article domain.Article, directory string) error {
	content, err := generateHTML(article, domain.ArticleHtmlTemplate)
	if err != nil {
		return fmt.Errorf("error generating html file content: %w", err)
	}

	if err := os.WriteFile(article.HtmlFilePath(directory), []byte(content), 0644); err != nil {
		return fmt.Errorf("error creating html file: %w", err)
	}
	return nil
}
