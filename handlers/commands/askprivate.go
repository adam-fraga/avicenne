package commands

//Command ask a question to be answer by the LLM and send an aaanswer in a private discord chan

import (
	"encoding/json"
	h "github.com/adam-fraga/avicenne/http"
	"github.com/bwmarrin/discordgo"
	"log"
)

func AskPrivate(discord *discordgo.Session, message *discordgo.MessageCreate, userPrompt string) error {
	// Send the request to the OpenAI API
	resultChan := make(chan string)
	errChan := make(chan error)

	go h.HttpRequestAsync(
		userPrompt,
		resultChan,
		errChan)

	channel, err := discord.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas pu créer un canal privé pour vous.")
		log.Printf("Erreur lors de la création d'un canal privé : %v", err)
		return err
	}
	select {
	case res := <-resultChan:
		var apiResponse map[string]interface{}
		if err := json.Unmarshal([]byte(res), &apiResponse); err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas pu traiter la réponse.")
			log.Printf("Erreur lors de l'analyse de la réponse API : %v", err)
			return err
		}
		content := apiResponse["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		discord.ChannelMessageSend(channel.ID, content)
	case err := <-errChan:
		discord.ChannelMessageSend(channel.ID, err.Error())
	}

	return nil
}
