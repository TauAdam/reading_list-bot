package main

import (
	"flag"
	telegramClient "github.com/tauadam/reading_list-bot/clients/telegram"
	event_consumer "github.com/tauadam/reading_list-bot/consumer/event-consumer"
	"github.com/tauadam/reading_list-bot/events/telegram"
	file_based "github.com/tauadam/reading_list-bot/storage/file-based"
	"log"
)

const (
	telegramApiHost = "https://api.telegram.org"
	PathToStorage   = "local-storage"
	BatchSize       = 100
)

func main() {
	tgClient := telegramClient.New(telegramApiHost, mustToken())

	localStorage := file_based.New(PathToStorage)

	eventsProcessor := telegram.New(tgClient, localStorage)

	log.Printf("program running")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
}

// mustToken returns a telegram bot token from the command line arguments
func mustToken() string {
	token := flag.String("tg-bot-token", "", "Telegram bot api token")

	if *token == "" {
		log.Fatalf("Token is required")
	}
	return *token
}
