package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"

	auth "github.com/adam-fraga/avicenne/handlers/auth"
	cmd "github.com/adam-fraga/avicenne/handlers/commands"
	llm "github.com/adam-fraga/avicenne/llm"
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
	llm.CurrentLLM.SetModel(os.Getenv("OPENAI_API_URL"), os.Getenv("GPT_TURBO"), os.Getenv("OPENAI_API_TOKEN"))
	fmt.Println("Set deffault LLM to chat gpt 3.5 turbo.")
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
	//WIPE MESSAGES
	case strings.Contains(userMessage, "!wipe"):
		isAdmin := auth.IsDiscordAdmin(discord, *message)
		log.Println("isAdmin: ", isAdmin)

		if !isAdmin {
			return
		}

		err := cmd.Wipe(discord, message)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			log.Println(err)
			return
		}
	//QUIT
	case strings.Contains(userMessage, "!quit"):
		isAdmin := auth.IsDiscordAdmin(discord, *message)

		if !isAdmin {
			return
		}

		// Send a confirmation message
		discord.ChannelMessageSend(message.ChannelID, "Avicen est maintenant hors ligne... Bye!")
		err := discord.Close()
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error while disconnecting: "+err.Error())
		}
		//SHOW ADMIN COMMANDA
	case strings.Contains(userMessage, "!admin"):
		isAdmin := auth.IsDiscordAdmin(discord, *message)

		if !isAdmin {
			return
		}
		cmd.HelpAdmin(discord, *message)

	//HELP
	case strings.Contains(userMessage, "!help"):
		cmd.Help(discord, *message)
	//ASK
	case strings.HasPrefix(userMessage, "!ask"):
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!ask"))
		err := cmd.Ask(discord, message, userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			return
		}
	//ASK PRIVATE
	case strings.Contains(userMessage, "!private"):
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!ask private"))
		err := cmd.AskPrivate(discord, message, userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			return
		}
	//TRANSLATE
	case strings.HasPrefix(userMessage, "!translate"):
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!translate"))
		err := cmd.Translate(discord, message, userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			return
		}
	//SPELLCHECK
	case strings.HasPrefix(userMessage, "!spellcheck"):
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!spellcheck"))
		err := cmd.Spellcheck(discord, message, userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			return
		}
	case strings.HasPrefix(userMessage, "!switchllm"):
		isAdmin := auth.IsDiscordAdmin(discord, *message)

		if !isAdmin {
			return
		}
		userPrompt := strings.TrimSpace(strings.TrimPrefix(userMessage, "!switchllm"))
		err := cmd.SwitchLLM(discord, message, userPrompt)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Désolé, une erreur est survenue. Réessaie plus tard.")
			return
		}
	case strings.HasPrefix(userMessage, "!showllm"):
		isAdmin := auth.IsDiscordAdmin(discord, *message)

		if !isAdmin {
			return
		}
		currentModel := llm.GetCurrentLLM()
		response := fmt.Sprintf("🤖 **LLM Actuel:** `%s`", currentModel.LLM)
		discord.ChannelMessageSend(message.ChannelID, response)
	}
}
