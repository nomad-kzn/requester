package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"requester/internal/entity"
	"strings"
)

func ParseCurlCmd(path string) (*entity.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	curlCmd := string(data)
	if !strings.HasPrefix(curlCmd, "curl") {
		return nil, errors.New("not a valid curl command")
	}

	reMethod := regexp.MustCompile(`-X\s+(\w+)`)
	reURL := regexp.MustCompile(`['"]?(https?://[^\s'"]+)['"]?`)
	reHeader := regexp.MustCompile(`-H\s+['"]?([^:'"]+):\s*([^'"]+)['"]?`)
	reData := regexp.MustCompile(`-d\s+['"](.+?)['"]`)

	method := http.MethodGet
	if m := reMethod.FindStringSubmatch(curlCmd); len(m) == 2 {
		method = m[1]
	}

	urlMatch := reURL.FindStringSubmatch(curlCmd)
	if len(urlMatch) != 2 {
		return nil, errors.New("URL not found")
	}
	fullURL := urlMatch[1]

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{}
	auth := (*entity.Authorization)(nil)
	for _, m := range reHeader.FindAllStringSubmatch(curlCmd, -1) {
		key := strings.TrimSpace(m[1])
		val := strings.TrimSpace(m[2])
		headers[key] = val
		if strings.ToLower(key) == "authorization" {
			parts := strings.SplitN(val, " ", 2)
			if len(parts) == 2 {
				auth = &entity.Authorization{Type: parts[0], Token: parts[1]}
			}
		}
	}

	body := map[string]any{}
	if m := reData.FindStringSubmatch(curlCmd); len(m) == 2 {
		if err := json.Unmarshal([]byte(m[1]), &body); err != nil {
			return nil, fmt.Errorf("failed to parse JSON body: %w", err)
		}
	}

	config := &entity.Config{
		RequestMethod:  method,
		URI:            parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path,
		RequestBody:    body,
		RequestHeaders: headers,
		Authorization:  auth,
	}

	if parsedURL.RawQuery != "" {
		q := parsedURL.RawQuery
		config.Query = &q
	}

	return config, nil
}
