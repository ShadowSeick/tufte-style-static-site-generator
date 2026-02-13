package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	gohttp "net/http"
)

func Get(ctx context.Context, req Request) ([]byte, error) {
	getReq, err := makeRequest(ctx, gohttp.MethodGet, req)
	if err != nil {
		return nil, fmt.Errorf("error creating get request: %w", err)
	}

	return sendRequest(getReq)
}

func Put(ctx context.Context, req Request) ([]byte, error) {
	fmt.Println(req.URL)
	putReq, err := makeRequest(ctx, gohttp.MethodPut, req)
	if err != nil {
		return nil, fmt.Errorf("error creating put request: %w", err)
	}

	return sendRequest(putReq)
}

func makeRequest(ctx context.Context, method string, request Request) (*gohttp.Request, error) {
	req, err := gohttp.NewRequestWithContext(ctx, method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func sendRequest(request *gohttp.Request) ([]byte, error) {
	response, err := gohttp.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
