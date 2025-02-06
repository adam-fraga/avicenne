package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	l "github.com/adam-fraga/avicenne/llm"
	"io"
	"net/http"
	"time"
)

func HttpRequestAsync(userMessage string, resultChan chan<- string, errChan chan<- error) {
	// Prepare HTTP request details
	httpMethod := "POST"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + l.CurrentLLM.ApiToken,
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"model": l.CurrentLLM.LLM,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": userMessage,
			},
		},
		"stream": false,
	}

	// Marshal the request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		errChan <- fmt.Errorf("error marshaling request body: %v", err)
		return
	}

	// Create the request
	req, err := http.NewRequest(httpMethod, l.CurrentLLM.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		errChan <- fmt.Errorf("error creating request: %v", err)
		return
	}

	// Set request headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request with a timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("error reading response body: %v", err)
		return
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("received non-OK HTTP status: %s, response: %s", resp.Status, string(responseBody))
		return
	}

	// Send the successful response back to the result channel
	resultChan <- string(responseBody)
}
