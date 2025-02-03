package auth

import "github.com/bwmarrin/discordgo"

func IsDiscordAdmin(discord *discordgo.Session, message discordgo.MessageCreate) bool {
	// Check if the user haas admin role
	for _, role := range message.Member.Roles {
		if role == "1333866427743998044" {
			return true
		}
	}

	discord.ChannelMessageSend(message.ChannelID, "Désolé, vous n'avez pas la permission d'effectuer cette action.")
	return false
}

func IsDiscordDeveloper(discord *discordgo.Session, message discordgo.MessageCreate) bool {
	// Check if the user haas developer role
	for _, role := range message.Member.Roles {
		if role == "1334949348781461584" {
			return true
		}
	}
	discord.ChannelMessageSend(message.ChannelID, "Désolé, vous n'avez pas la permission d'effectuer cette action.")
	return false
}
