package file_based

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tauadam/reading_list-bot/lib/utils"
	"github.com/tauadam/reading_list-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
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

	dirPath := filepath.Join(s.basePath, article.UserName)

	if err := os.MkdirAll(dirPath, 0744); err != nil {
		return err
	}

	fileName, err := generateFileName(article)
	if err != nil {
		return err
	}

	pathToFile := filepath.Join(dirPath, fileName)

	file, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(article); err != nil {
		return err
	}

	return nil
}

func generateFileName(a *storage.Article) (string, error) {
	return a.Hash()
}

var ErrNotFound = errors.New("no saved articles")

func (s Storage) PickRandom(userName string) (article *storage.Article, err error) {
	defer func() { err = utils.Wrap("can't pick random article", err) }()

	dirPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNotFound
	}

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodeArticle(filepath.Join(dirPath, file.Name()))
}

func (s Storage) decodeArticle(pathToFile string) (*storage.Article, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, utils.Wrap("can't open file", err)
	}

	defer func() { _ = file.Close() }()

	var article storage.Article
	if err := gob.NewDecoder(file).Decode(&article); err != nil {
		return nil, utils.Wrap("can't decode article", err)
	}

	return &article, nil
}

func (s Storage) Remove(a *storage.Article) error {
	fileName, err := generateFileName(a)
	if err != nil {
		return utils.Wrap("can't remove", err)
	}

	pathToFile := filepath.Join(s.basePath, a.UserName, fileName)

	if err := os.Remove(pathToFile); err != nil {
		return utils.Wrap(fmt.Sprintf("can't remove file %s", pathToFile), err)
	}

	return nil
}
