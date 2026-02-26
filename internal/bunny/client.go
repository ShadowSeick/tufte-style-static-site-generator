package bunny

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ShadowSeick/tufte-style-static-site-generator/internal/http"
)

// Implementing https://docs.bunny.net/api-reference/storage
const (
	AuthorizationHeader = "AccessKey"
	ChecksumHeader      = "Checksum"
)

var (
	initialize      sync.Once
	bunnyStorageURL string
	headers         = make(map[string]string, 2) // We will have a maximum of 2 headers, access key and checksum
)

func Init() {
	initialize.Do(func() {
		bunnyStorageURL = os.Getenv("BUNNY_STORAGE_API_URL")
		if bunnyStorageURL == "" {
			panic("bunny storage api url not set")
		}

		storageZone := os.Getenv("BUNNY_STORAGE_ZONE")
		if storageZone == "" {
			panic("bunny storage zone not set")
		}
		bunnyStorageURL = fmt.Sprintf("%s%s/", bunnyStorageURL, storageZone)

		accessKey := os.Getenv("BUNNY_STORAGE_ACCESS_KEY")
		if accessKey == "" {
			panic("bunny storage access key not set")
		}
		headers[AuthorizationHeader] = accessKey
	})
}

func GetFile(ctx context.Context, path string) ([]File, error) {
	body, err := http.Get(ctx, http.Request{
		URL:     bunnyStorageURL + path + "/",
		Headers: headers,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting file from bunny Storage: %w", err)
	}

	var files []File
	if err := json.Unmarshal(body, &files); err != nil {
		var httpError HTTPError
		if subErr := json.Unmarshal(body, &httpError); subErr != nil {
			return nil, fmt.Errorf("error unmarshalling files: %w", err)
		}
		if !httpError.NotFound() {
			return nil, fmt.Errorf("error sending request:\n httpCode: %d\n message: %s", httpError.Code, httpError.Message)
		}
	}
	return files, nil
}

func UploadFile(ctx context.Context, path string, checksum string, data []byte) error {
	uploadHeaders := make(map[string]string, 2)
	uploadHeaders[AuthorizationHeader] = headers[AuthorizationHeader]
	uploadHeaders[ChecksumHeader] = checksum
	uploadHeaders[http.ContentType.String()] = http.ApplicationOctetStream.String()

	body, err := http.Put(ctx, http.Request{
		URL:     bunnyStorageURL + path,
		Headers: uploadHeaders,
		Body:    data,
	})
	if err != nil {
		var httpError HTTPError
		if subErr := json.Unmarshal(body, &httpError); subErr != nil {
			return fmt.Errorf("error unmarshalling files: %w", err)
		}
		return fmt.Errorf("error sending request:\n httpCode: %d\n message: %s", httpError.Code, httpError.Message)
	}
	fmt.Println(string(body))
	return nil
}
