package main

import (
	"fmt"
	"context"
	"sync"
	"time"

	"github.com/ShadowSeick/tufte-style-static-site-generator/domain"
	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/bunny"
)

func InitUpload(ctx context.Context, files []domain.File) {
	bunny.Init()

	newCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Go(func() {
			files, err := bunny.GetFile(newCtx, file.RemoteFilePath())
			if err != nil {
				fmt.Println("error while getting file: ", err)
				return
			}

			if len(files) > 1 {
				fmt.Println("more than one file with the same name! ", files)
				return
			}

			if len(files) == 1 && files[0].Checksum != nil && *files[0].Checksum == file.Checksum() {
				fmt.Println("file ", file.Title(), " is already uploaded")
				return
			}

			fmt.Println("uploading file...", file.Title())
			if err := bunny.UploadFile(ctx, file.RemoteFilePath(), file.Checksum(), []byte(file.Content)); err != nil {
				fmt.Println("error updating file: %w", err)
				return
			}
		})
	}

	wg.Wait()

	fmt.Println("Finished uploading articles")
}
