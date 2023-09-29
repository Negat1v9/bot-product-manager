package main

import (
	"log"
	"os"

	telegram "github.com/Negat1v9/telegram-bot-orders/internal/bot"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	timeOut = 5
	offset  = 0
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Env files is not found: %s", err.Error())
	}

}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	dbURL := os.Getenv("DB_URL")

	client, err := tgbot.NewBotAPI(token)

	if err != nil {
		log.Fatalf("Can't create telegram bot Client: %s", err.Error())
	}

	bot := telegram.New(client, timeOut, offset)

	log.Println("Start pooling Bot")

	if err = bot.Start(dbURL); err != nil {
		log.Fatalln(err.Error())
	}
}
