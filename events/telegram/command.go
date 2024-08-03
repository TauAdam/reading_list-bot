package telegram

import (
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
