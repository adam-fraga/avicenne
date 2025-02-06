package commands

import (
	"fmt"
	"os"
	"strings"

	llm "github.com/adam-fraga/avicenne/llm"
	"github.com/bwmarrin/discordgo"
)

func SwitchLLM(discord *discordgo.Session, msg *discordgo.MessageCreate, userPrompt string) error {
	parts := strings.SplitN(userPrompt, " ", 2)
	model := parts[0]

	if len(parts) < 1 {
		discord.ChannelMessageSend(msg.ChannelID, "❌ Veuillez fournir le nom du modèle que vous voulez utiliser.")
		return fmt.Errorf("invalid format for switch LLM command")
	}

	switch model {
	case "gpt-3.5":
		llm.CurrentLLM.SetModel(os.Getenv("OPENAI_API_URL"), os.Getenv("GPT_TURBO"), os.Getenv("OPENAI_API_TOKEN"))
	case "gpt-4":
		llm.CurrentLLM.SetModel(os.Getenv("OPENAI_API_URL"), os.Getenv("GPT4"), os.Getenv("OPENAI_API_TOKEN"))
	case "deepseek-v3":
		llm.CurrentLLM.SetModel(os.Getenv("DS_API_URL"), os.Getenv("DSV3"), os.Getenv("DS_API_TOKEN"))
	case "deepseek-r1":
		llm.CurrentLLM.SetModel(os.Getenv("DS_API_URL"), os.Getenv("DSR1"), os.Getenv("DS_API_TOKEN"))
	case "sonnet-3.5":
		llm.CurrentLLM.SetModel(os.Getenv("CLAUDE_API_URL"), os.Getenv("SONET"), os.Getenv("CLAUDE_API_TOKEN"))
	default:
		discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("❌ Modèle: %s non supporté.", model))
		return fmt.Errorf("Error Model %s does not exist !", model)
	}

	discord.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("✅ Modèle LLM changé en **%s**", llm.CurrentLLM.LLM))
	return nil
}
