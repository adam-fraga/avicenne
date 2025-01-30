package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
)

func Translate(discord *discordgo.Session, message *discordgo.MessageCreate, userPrompt string) error {
	parts := strings.SplitN(userPrompt, " ", 2) // Split into 2 parts: language and text
	if len(parts) < 2 {
		discord.ChannelMessageSend(message.ChannelID, "Veuillez fournir un texte à traduire.")
		return fmt.Errorf("invalid format for translation command")
	}
	resultChan := make(chan string)
	errChan := make(chan error)

	targetLanguage := parts[0]
	textToTranslate := parts[1] // The actual text to translate

	translationPrompt := fmt.Sprintf("Translate the following text into %s: %s", targetLanguage, textToTranslate)

	discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Je vais traduire ça en %s pour toi...", targetLanguage))

	go AskHttpRequestAsync(os.Getenv("API_URL"), os.Getenv("API_TOKEN"), translationPrompt, resultChan, errChan)

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
