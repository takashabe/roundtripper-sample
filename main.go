package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func main() {
	if err := runSimple(); err != nil {
		panic(err)
	}

	if err := runLogging(); err != nil {
		panic(err)
	}
}

func runLogging() error {
	client := &http.Client{
		Transport: &LoggingRoundTripper{
			MaskFields: []string{"password"},
		},
	}
	body := map[string]string{
		"username": "user",
		"password": "password",
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://httpbin.org/post", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func runSimple() error {
	client := &http.Client{
		Transport: &SimpleRoundTripper{},
	}
	_, err := client.Get("http://example.com")
	if err != nil {
		return err
	}
	return nil
}
