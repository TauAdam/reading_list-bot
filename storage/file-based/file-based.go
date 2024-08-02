package file_based

import (
	"github.com/tauadam/reading_list-bot/lib/utils"
	"github.com/tauadam/reading_list-bot/storage"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{
		basePath,
	}
}

func (s Storage) Save(article *storage.Article) (err error) {
	defer func() { err = utils.Wrap("can't save", err) }()

	filePath := filepath.Join(s.basePath, article.UserName)

	if err := os.MkdirAll(filePath, 0744); err != nil {
		return err
	}

	fileName, err := generateFileName(article)
	if err != nil {
		return err
	}

	//	TODO write to file
}

func generateFileName(a *storage.Article) (string, error) {
	return a.Hash()
}
