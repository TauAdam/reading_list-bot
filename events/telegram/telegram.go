package telegram

import (
	"errors"
	"github.com/tauadam/reading_list-bot/clients/telegram"
	"github.com/tauadam/reading_list-bot/events"
	"github.com/tauadam/reading_list-bot/lib/utils"
	"github.com/tauadam/reading_list-bot/storage"
)

type Meta struct {
	ChatID   int
	UserName string
}

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
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, utils.Wrap("can't fetch updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, updateToEvent(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func updateToEvent(update telegram.Update) events.Event {
	updateType := fetchType(update)

	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return res
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

// Custom error types for events like messages
var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown extractMeta type")
)

// Process realization that works with Telegram API, processes messages
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return utils.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := extractMeta(event)
	if err != nil {
		return utils.Wrap("can't extract meta info", err)
	}

	if err := p.cmdRouter(event.Text, meta.ChatID, meta.UserName); err != nil {
		return utils.Wrap("can't execute command", err)
	}

	return nil
}

// extractMeta extracts chat id and user name from api
func extractMeta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, utils.Wrap("can't extract meta info", ErrUnknownMetaType)
	}

	return res, nil
}
