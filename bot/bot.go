package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Empêcher le bot de répondre à ses propres messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Extraire le contenu du message de l'utilisateur
	userMessage := message.Content

	// Gérer les différentes commandes
	switch {
	case strings.Contains(userMessage, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Salut ! Je suis Avicenne, ton assistant intelligent 🤖. Besoin d'aide ? Tape !help 🚀. Une question ? Utilise !ask 💡")
	case strings.HasPrefix(userMessage, "!ask"):
		// Supprimer le préfixe "!ask " du message de l'utilisateur
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!ask"))

		// Envoyer la requête à l'API OpenAI
		res, err := SendRequest(os.Getenv("API_URL"), os.Getenv("API_TOKEN"), userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			log.Printf("Erreur lors de l'envoi de la requête : %v", err)
			return
		}

		// Analyser la réponse de l'API
		var apiResponse map[string]interface{}
		if err := json.Unmarshal([]byte(res), &apiResponse); err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas pu traiter la réponse.")
			log.Printf("Erreur lors de l'analyse de la réponse API : %v", err)
			return
		}

		// Extraire le message de l'assistant de la réponse
		choices, ok := apiResponse["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas obtenu de réponse valide.")
			log.Printf("Réponse API invalide : %v", apiResponse)
			return
		}

		firstChoice, ok := choices[0].(map[string]interface{})
		if !ok {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas pu interpréter la réponse.")
			log.Printf("Format de choix invalide : %v", choices[0])
			return
		}

		assistantMessage, ok := firstChoice["message"].(map[string]interface{})
		if !ok {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, je n'ai pas trouvé la réponse de l'assistant.")
			log.Printf("Format du message invalide : %v", firstChoice["message"])
			return
		}

		content, ok := assistantMessage["content"].(string)
		if !ok {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, la réponse de l'assistant est invalide.")
			log.Printf("Format du contenu invalide : %v", assistantMessage["content"])
			return
		}

		// Send the assistant's message to the Discord channel
		response := fmt.Sprintf("**🗣️ Ta question :** %s\n\n%s", userPrompt, content)
		discord.ChannelMessageSend(message.ChannelID, response)
	}
}
