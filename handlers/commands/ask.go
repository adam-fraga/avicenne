package commands

//Command that ask a regular question and use the LLM to answer

import (
	"encoding/json"
	"log"

	h "github.com/adam-fraga/avicenne/http"
	"github.com/bwmarrin/discordgo"
)

//Command that switch LLM

func Ask(discord *discordgo.Session, message *discordgo.MessageCreate, userPrompt string) error {
	resultChan := make(chan string)
	errChan := make(chan error)

	go h.HttpRequestAsync(
		userPrompt,
		resultChan,
		errChan)

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
