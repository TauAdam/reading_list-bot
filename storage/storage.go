package storage

type Storage interface {
	Save(a *Article) error
	PickRandom(userName string) (*Article, error)
	Remove()
	IsExist() bool
}

type Article struct {
	URL      string
	UserName string
}
