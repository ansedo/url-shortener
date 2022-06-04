package shortener

import (
	"github.com/ansedo/url-shortener/internal/storage"
	"strconv"
)

type Shortener struct {
	Storage storage.Storage
}

func NewShortener(storage storage.Storage) *Shortener {
	return &Shortener{
		Storage: storage,
	}
}

func (s *Shortener) GenerateID() (string, error) {
	// Strong and tiny approach with extra short ids (at least for first 10 values) and no collisions!
	return strconv.Itoa(s.Storage.Count()), nil
}
