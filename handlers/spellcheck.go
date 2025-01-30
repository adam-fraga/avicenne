package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func Spellcheck(discord *discordgo.Session, message *discordgo.MessageCreate, userPrompt string) error {
	resultChan := make(chan string)
	errChan := make(chan error)

	spellcheckPrompt := fmt.Sprintf("Please check the spelling and grammar of the following text:\n\n %s", userPrompt)
	discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Je vais vérifier l'orthographe et la grammaire du texte suivant pour toi..."))

	go AskHttpRequestAsync(os.Getenv("API_URL"), os.Getenv("API_TOKEN"), spellcheckPrompt, resultChan, errChan)

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
