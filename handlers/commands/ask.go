package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func AskHttpRequestAsync(url string, apiKey string, userMessage string, resultChan chan<- string, errChan chan<- error) {
	// Prepare HTTP request details
	httpMethod := "POST"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
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
	}

	// Marshal the request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		errChan <- fmt.Errorf("error marshaling request body: %v", err)
		return
	}

	// Create the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBody))
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
		Timeout: 10 * time.Second, // 10 seconds timeout
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

func Ask(discord *discordgo.Session, message *discordgo.MessageCreate, userPrompt string) error {
	resultChan := make(chan string)
	errChan := make(chan error)

	go AskHttpRequestAsync(os.Getenv("API_URL"), os.Getenv("API_TOKEN"), userPrompt, resultChan, errChan)

	select {
	case res := <-resultChan:
		var apiResponse map[string]interface{}
		if err := json.Unmarshal([]byte(res), &apiResponse); err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas pu traiter la réponse.")
			log.Printf("Erreur lors de l'analyse de la réponse API : %v", err)
			return err
		}
		content := apiResponse["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		discord.ChannelMessageSend(message.ChannelID, content)
	case err := <-errChan:
		discord.ChannelMessageSend(message.ChannelID, err.Error())
	}

	return nil
}
