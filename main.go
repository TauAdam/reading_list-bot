package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()

}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "Telegram bot api token")

	if *token == "" {
		log.Fatalf("Token is required")
	}
	return *token
}
