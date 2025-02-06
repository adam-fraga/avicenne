package commands

//Display Help message on how to use Avicen

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Help(discord *discordgo.Session, message discordgo.MessageCreate) {

	discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		`Salut %s! Je suis Avicen ton bot d'assistance intelligent.
  
  Voici ce que je peux faire pour toi :

  -🚀 Besoin d'aide ? Tape !help.

  -💡 Une question ? Tape !ask "Ta question".

  -📬 Tu souhaites poser ta question en privé ? Tape !private "Ta question".

  -🌍 Besoin de traduction ? Tape !translate langue "texte à traduire".

  --✍️  Besoin d'aide pour la correction d'orthographe ? Tape !spellcheck "texte à corriger".
  `, message.Author.Username))
}

func HelpAdmin(discord *discordgo.Session, message discordgo.MessageCreate) {
	discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(
		`Salut %s ! Voici la liste des commandes Admin.
  
  Voici ce que je peux faire pour toi :

  -🚀 Besoin d'aide ? Tape !admin.

  -📂 Enregistrer un document sur google drive ? Tape !store "documentname".  

  -🔄 Changer de modèle IA ? Tape !switchllm "nom du modele" (gpt-3.5, gpt-4, deepseek-r1, deepseek-v3, sonet-3.5)

  -🔄 Voir le LLM courant ? Tape !showllm
  `, message.Author.Username))
}
