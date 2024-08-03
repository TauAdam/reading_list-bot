package telegram

import (
	"github.com/tauadam/reading_list-bot/clients/telegram"
	"github.com/tauadam/reading_list-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}
