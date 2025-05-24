package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	Config struct {
		RequestMethod  string            `json:"request_method"`
		URI            string            `json:"uri"`
		Query          *string           `json:"query,omitempty"`
		RequestBody    map[string]any    `json:"request_body"`
		RequestHeaders map[string]string `json:"request_headers"`
		Authorization  *Authorization    `json:"authorization,omitempty"`
	}

	Authorization struct {
		Type  string
		Token string
	}
)

func (c *Config) MakeRequestURI() string {
	if c.Query != nil {
		return fmt.Sprintf("%s?%s", c.URI, *c.Query)
	}

	return c.URI
}

func (c *Config) MakeRequestBody() (io.Reader, error) {
	if c.RequestBody == nil {
		return nil, nil
	}

	reqBodyBytes, err := json.Marshal(c.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return io.NopCloser(bytes.NewReader(reqBodyBytes)), nil
}

func (c *Config) AddHeadersToRequest(r *http.Request) {
	for k, v := range c.RequestHeaders {
		r.Header.Add(k, v)
	}
}
