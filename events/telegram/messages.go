package telegram

const msgHelp = `I can save and keep your articles. Also I can offer you them to read.

In order to save the article, just send me link to it. After that article will be saved in database.

In order to get a random article from your list, send me "/rand" command

Caution! After that, that article will be removed`

const msgGreeting = "Hi there!\n\n I'm reading list helper bot.\n\n" + msgHelp

const (
	msgUnknown       = "Unknown command 🤔"
	msgNotExists     = "No articles saved 🧐"
	msgSaved         = "Saved 👌🏾"
	msgAlreadyExists = "This article already saved 🤗"
)
