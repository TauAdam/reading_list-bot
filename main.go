package main

import (
	"flag"
	"github.com/tauadam/reading_list-bot/clients/telegram"
	"log"
)

func main() {
	token := mustToken()

	tgClient := telegram.New("https://api.telegram.org", token)
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "Telegram bot api token")

	if *token == "" {
		log.Fatalf("Token is required")
	}
	return *token
}
