package main

import (
	"context"
	telegramClient "github.com/tauadam/reading_list-bot/clients/telegram"
	event_consumer "github.com/tauadam/reading_list-bot/consumer/event-consumer"
	"github.com/tauadam/reading_list-bot/events/telegram"
	"github.com/tauadam/reading_list-bot/storage/sqlite"
	"log"
	"os"
)

const (
	telegramApiHost = "api.telegram.org"
	PathToStorage   = "local-storage"
	PathToSqliteDb  = "database/sqlite.db"
	BatchSize       = 100
)

func main() {
	tgClient := telegramClient.New(telegramApiHost, mustToken())

	//localStorage := file_based.New(PathToStorage)

	sqliteStorage, err := sqlite.New(PathToSqliteDb)
	if err != nil {
		log.Fatalf("failed to create sqlite storage: %v", err)
	}

	if err := sqliteStorage.Init(context.TODO()); err != nil {
		log.Fatalf("failed to init sqlite storage: %v", err)
	}

	eventsProcessor := telegram.New(tgClient, sqliteStorage)

	log.Printf("program running")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
}

// mustToken returns a telegram bot token from the command line arguments
func mustToken() string {
	token := os.Getenv("TG_BOT_TOKEN")

	if token == "" {
		log.Fatalf("Token is required")
	}

	return token
}
