package main

import (
	bot "github.com/adam-fraga/avicenne/bot"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot.BotToken = os.Getenv("BOT_TOKEN")
	// apiToken := os.Getenv("API_TOKEN")

	bot.Run() // call the run function of bot/bot.go
}
