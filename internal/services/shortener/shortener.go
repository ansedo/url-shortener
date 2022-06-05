package shortener

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/storage"
	"strconv"
)

type Shortener struct {
	Config  *config.Config
	Storage storage.Storager
}

func NewShortener(opts ...Option) *Shortener {
	s := &Shortener{}

	for _, opt := range opts {
		opt(s)
	}

	if s.Config == nil {
		WithDefaultConfig()(s)
	}

	if s.Storage == nil {
		WithDefaultStorage()(s)
	}

	return s
}

func (s *Shortener) GenerateID() (string, error) {
	// Strong and tiny approach with extra short ids (at least for first 10 values) and no collisions!
	return strconv.Itoa(s.Storage.Count()), nil
}
