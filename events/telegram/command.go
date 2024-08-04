package telegram

import (
	"context"
	"errors"
	"github.com/tauadam/reading_list-bot/clients/telegram"
	"github.com/tauadam/reading_list-bot/lib/utils"
	"github.com/tauadam/reading_list-bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RandomArticleCmd = "/rand"
	HelpCmd          = "/help"
	StartCmd         = "/start"
)

func (p *Processor) cmdRouter(text string, chatID int, userName string) error {
	command := strings.TrimSpace(text)

	log.Printf("got new command: [%s] from %s", command, userName)

	if isValidAddCmd(command) {
		return p.handleSave(chatID, userName, command)
	}

	switch command {
	case StartCmd:
		return p.handleStart(chatID)
	case RandomArticleCmd:
		return p.handleRandom(chatID, userName)
	case HelpCmd:
		return p.handleHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknown)
	}
}

// isValidAddCmd checks if user typed valid command to add new article
func isValidAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

// handleSave saves data and send user status about operation
func (p *Processor) handleSave(chatID int, articleURL string, userName string) (err error) {
	defer func() { err = utils.Wrap("can't execute save operation", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	article := &storage.Article{
		URL:      articleURL,
		UserName: userName,
	}

	isExists, err := p.storage.IsExist(context.Background(), article)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(context.Background(), article); err != nil {
		return err
	}

	if err := sendMsg(msgSuccessfullySaved); err != nil {
		return err
	}

	return nil
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(message string) error {
		return tg.SendMessage(chatID, message)
	}
}

func (p *Processor) handleRandom(chatID int, userName string) (err error) {
	defer func() { err = utils.Wrap("can't execute random operation", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	article, err := p.storage.PickRandom(context.Background(), userName)
	if err != nil && !errors.Is(err, storage.ErrArticleNotFound) {
		return err
	}

	if errors.Is(err, storage.ErrArticleNotFound) {
		return sendMsg(msgNotExists)
	}

	if err = sendMsg(article.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), article)
}

func (p *Processor) handleHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) handleStart(chatID int) error {
	return p.tg.SendMessage(chatID, msgGreeting)
}
