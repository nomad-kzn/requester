package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
)

type (
	HTTPResponse struct {
		StatusCode    int
		Headers       map[string]string
		Body          []byte
		ContentLength int
		Duration      time.Duration
		URL           string
		Method        string
	}
)

func MakeHTTPResponse(resp *http.Response, timeStart time.Time) (*HTTPResponse, error) {
	if resp == nil {
		return nil, nil
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &HTTPResponse{
		StatusCode:    resp.StatusCode,
		Headers:       headers,
		Body:          body,
		ContentLength: len(body),
		Duration:      time.Since(timeStart),
		URL:           resp.Request.URL.String(),
		Method:        resp.Request.Method,
	}, nil
}

func (r *HTTPResponse) PrintSummary() {
	fmt.Println(colorGreen, "===== HTTP Response =====", colorReset)

	fmt.Printf("%sMethod:%s        %s\n", colorCyan, colorReset, r.Method)
	fmt.Printf("%sURL:%s           %s\n", colorCyan, colorReset, r.URL)

	statusColor := colorGreen
	if r.StatusCode >= 400 && r.StatusCode < 500 {
		statusColor = colorYellow
	} else if r.StatusCode >= 500 {
		statusColor = colorRed
	}

	fmt.Printf("%sStatus Code:%s   %s%d%s\n", colorCyan, colorReset, statusColor, r.StatusCode, colorReset)
	fmt.Printf("%sDuration:%s      %s\n", colorCyan, colorReset, r.Duration)
	fmt.Printf("%sSize:%s          %d bytes\n", colorCyan, colorReset, r.ContentLength)

	if len(r.Headers) > 0 {
		fmt.Printf("%sHeaders:%s\n", colorCyan, colorReset)
		for k, v := range r.Headers {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	if len(r.Body) > 0 {
		contentType := r.Headers["Content-Type"]
		fmt.Printf("%sBody (preview):%s\n", colorCyan, colorReset)

		if strings.Contains(contentType, "application/json") {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, r.Body, "", "  ")
			if err != nil {
				// если не получилось красиво, просто выводим
				fmt.Println(string(r.Body))
			} else {
				fmt.Println(prettyJSON.String())
			}
		} else {
			// Обычный текст
			preview := string(r.Body)
			if len(preview) > 500 {
				preview = preview[:500] + "..."
			}
			fmt.Println(strings.TrimSpace(preview))
		}
	}
	fmt.Println()
}
