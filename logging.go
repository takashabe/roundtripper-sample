package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type LoggingRoundTripper struct {
	MaskFields []string
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	clone := req.Clone(ctx)
	b, err := io.ReadAll(clone.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer clone.Body.Close()

	// 元リクエストのBodyをリストア
	req.Body = io.NopCloser(bytes.NewReader(b))

	// dump request
	maskedBody, err := maskJSONFields(b, lrt.MaskFields)
	if err != nil {
		return nil, fmt.Errorf("failed to mask request body: %w", err)
	}
	clone.Body = io.NopCloser(bytes.NewReader(maskedBody))
	clone.ContentLength = int64(len(maskedBody))
	clone.Header.Set("Content-Length", fmt.Sprint(len(maskedBody)))
	reqDump, err := httputil.DumpRequestOut(clone, true)
	if err != nil {
		return nil, fmt.Errorf("httputil.DumpRequestOut: %w", err)
	}
	log.Println(string(reqDump))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultTransport.RoundTrip: %w", err)
	}

	// dump response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, fmt.Errorf("httputil.DumpResponse: %w", err)
	}
	log.Println(string(respDump))

	return resp, nil
}

// maskJSONFields request bodyの指定されたフィールドをマスクする. bodyがJSONでない場合は何もしない
func maskJSONFields(body []byte, fieldsToMask []string) ([]byte, error) {
	if ok := json.Valid(body); !ok {
		return body, nil
	}

	maskedBody := string(body)
	var err error
	for _, field := range fieldsToMask {
		if !gjson.Get(maskedBody, field).Exists() {
			continue
		}
		maskedBody, err = sjson.Set(maskedBody, field, "***")
		if err != nil {
			return nil, fmt.Errorf("failed to mask field: %w", err)
		}
	}
	return []byte(maskedBody), nil
}
