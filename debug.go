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
	"github.com/ShadowSeick/tufte-style-static-site-generator/generate"
)

var (
	changesChannel = make(chan int)
	initialize sync.Once
	debugDirectories = []string{
		"debug/en/articles",
		"debug/es/articles",
	}
)

func InitDebug(ctx context.Context, files []domain.File) {
	initialize.Do(func() {
		// Create directories for files
		cwd, _ := os.Getwd()
		for _, dir := range debugDirectories {
			if err := os.MkdirAll(filepath.Join(cwd, dir), 0755); err != nil && !os.IsExist(err) {
				log.Fatalf("error creating directory %s: %v", dir, err)
			}
		}

		newCtx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		for index, file := range files {
			if err := os.WriteFile(file.DebugFilePath(), []byte(file.Content), 0644); err != nil {
				log.Fatalf("error creating html file: %v", err)
			}

			// Only hot reload articles
			if file.Type != domain.Article {
				continue
			}

			wg.Go(func () {
				defer wg.Done()
				listenToArticleChanges(newCtx, index, file, cancel)
			})

			wg.Go(func () {
				defer close(changesChannel)
				wg.Wait()
			})
		}

		go regenerateHTMLFiles(ctx, files)

		// Run server
		mux := http.NewServeMux()

		page := http.FileServer(http.Dir("debug"))
		mux.Handle("/", page)

		assets := http.FileServer(http.Dir("assets"))
		mux.Handle("/assets/", http.StripPrefix("/assets/", assets))
		
		log.Print("Listening on :3000...")
		log.Fatal(http.ListenAndServe(":3000", mux))
	})
}

func regenerateHTMLFiles(ctx context.Context, files []domain.File) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("gracefully shutting down")
			return
		case articleIndex, ok := <-changesChannel:
			if !ok { // Channel closed
				return
			}

			content, err := generate.ArticleHtml(files[articleIndex])
			if err != nil {
				fmt.Println("error regenerating html file", files[articleIndex].Name, err)
				return
			}
			files[articleIndex].Content = content

			if err := os.WriteFile(files[articleIndex].DebugFilePath(), []byte(files[articleIndex].Content), 0644); err != nil {
				log.Fatalf("error creating html file: %v", err)
			}
		}
	}
}

func listenToArticleChanges(ctx context.Context, index int, file domain.File, cancel context.CancelFunc) {
	filepath := file.LocalFilePath()
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
			fmt.Println("gracefully shuttinwdown listening to article: ", file.Title())
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
				changesChannel <- index
			}
		}
	}
}
