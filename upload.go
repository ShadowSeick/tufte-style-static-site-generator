package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/bunny"
)

func InitUpload(ctx context.Context, files []domain.File) {
	bunny.Init()

	newCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Get the remote files and do your fucking job
	uploadedFiles := make(map[string]bunny.File, 0)
	availableFiles := []domain.FileType{domain.Article, domain.Index, domain.ArticleIndex, domain.ContactIndex}
	// Surely there is a better way, but it's not necessary for now
	// For each language get all remote files
	for i := range domain.LanguageCount {
		for _, fileType := range availableFiles {
			files, err := bunny.GetFile(newCtx, filepath.Join(i.String(), fileType.RemoteDirectory()))
			if err != nil {
				fmt.Println("error while getting file: ", err)
				return
			}

			for _, file := range files {
				if file.IsDirectory {
					continue
				}
				uploadedFiles[filepath.Join(i.String(), fileType.RemoteDirectory(), file.Name)] = file
			}
		}
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Go(func() {
			uploadedFile := uploadedFiles[file.RemoteFilePath()]
			if uploadedFile.Checksum != nil && *uploadedFile.Checksum == file.Checksum() {
				fmt.Println("file ", file.Title(), " is already uploaded")
				return
			}

			fmt.Println("uploading file... ", file.Title())
			if err := bunny.UploadFile(ctx, file.RemoteFilePath(), file.Checksum(), []byte(file.Content)); err != nil {
				fmt.Println("error updating file: %w", err)
				return
			}
		})
	}

	wg.Wait()

	fmt.Println("Finished uploading files")
}
