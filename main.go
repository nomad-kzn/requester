package main

import (
	"context"
	"log"
	"net/http"
	entity2 "requester/internal/entity"
	"requester/internal/usecase"
	"time"
)

func main() {
	ctx := context.Background()

	config, err := usecase.ParseCurlCmd("./req.curl")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
		return
	}

	reqBody, err := config.MakeRequestBody()
	if err != nil {
		log.Fatalf("failed to make request body: %v", err)
		return
	}

	request, err := http.NewRequestWithContext(ctx, config.RequestMethod, config.MakeRequestURI(), reqBody)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
		return
	}
	config.AddHeadersToRequest(request)

	startTime := time.Now()
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("failed to make request: %v", err)
		return
	}
	defer resp.Body.Close()

	httpResponse, err := entity2.MakeHTTPResponse(resp, startTime)
	if err != nil {
		log.Fatalf("failed to make http response: %v", err)
		return
	}

	httpResponse.PrintSummary()
}
