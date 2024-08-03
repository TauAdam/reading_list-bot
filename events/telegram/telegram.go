package telegram

import (
	"github.com/tauadam/reading_list-bot/clients/telegram"
	"github.com/tauadam/reading_list-bot/events"
	"github.com/tauadam/reading_list-bot/lib/utils"
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

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, utils.Wrap("can't fetch updates", err)
	}

	res := make([]events.Event, 0, len(update))

}

func updateToEvent(update telegram.Update) events.Event {
	return events.Event{
		Type: fetchType(update),
		Text: fetchText(update),
	}
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}
