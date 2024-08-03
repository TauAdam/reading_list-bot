package telegram

import (
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

func (p *Processor) cmdRouter(command string, chatID int, userName string) error {
	command := strings.TrimSpace(command)

	log.Printf("got new command: [%s] from %s", command, userName)

	if isValidAddCmd(command) {
		//	TODO: add new article
	}
	switch command {
	case StartCmd:

	case RandomArticleCmd:

	case HelpCmd:

	default:

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

	isExists, err := p.storage.IsExist(article)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(article); err != nil {
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

	article, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrArticleNotFound) {
		return err
	}

	if errors.Is(err, storage.ErrArticleNotFound) {
		return sendMsg(msgNotExists)
	}

	if err = sendMsg(article.URL); err != nil {
		return err
	}

	return p.storage.Remove(article)
}
