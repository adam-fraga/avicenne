package handlers

import "github.com/bwmarrin/discordgo"

func Help(discord *discordgo.Session, message discordgo.MessageCreate) {

	discord.ChannelMessageSend(message.ChannelID,
		`Salut ! Je suis Avicenne, ton assistant intelligent 🤖.
  
  Voici ce que je peux faire pour toi :

  -🚀 Besoin d'aide ? Tape !help.

  -💡 Une question ? Tape !ask "Ta question".

  -📬 Tu souhaites poser ta question en privé ? Tape !private "Ta question".

  -🌍 Besoin de traduction ? Tape !translate langue "texte à traduire".

  -✍️  Besoin d'aide pour la correction d'orthographe ? Tape !spellcheck "texte à corriger".
  `)
}
