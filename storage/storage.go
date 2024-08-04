package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/tauadam/reading_list-bot/lib/utils"
	"io"
)

var ErrArticleNotFound = errors.New("no saved articles")

type Storage interface {
	Save(a *Article) error
	PickRandom(userName string) (*Article, error)
	Remove(a *Article) error
	IsExist(a *Article) (bool, error)
}

type Article struct {
	URL      string
	UserName string
}

// Hash returns a hash of the article
func (a Article) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, a.URL); err != nil {
		return "", utils.Wrap("can't calc hash", err)
	}

	if _, err := io.WriteString(h, a.UserName); err != nil {
		return "", utils.Wrap("can't calc hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
