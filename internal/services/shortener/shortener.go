package shortener

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/storages"
	"strconv"
)

type Shortener struct {
	Storage storages.Storager
}

func New(opts ...Option) *Shortener {
	s := &Shortener{}

	for _, opt := range opts {
		opt(s)
	}

	if s.Storage == nil && config.New().FileStoragePath != "" {
		WithFileStorage()(s)
	}

	if s.Storage == nil {
		WithDefaultStorage()(s)
	}

	return s
}

func (s *Shortener) GenerateID() (string, error) {
	// Strong and tiny approach with extra short ids (at least for first 10 values) and no collisions!
	return strconv.Itoa(s.Storage.NextID()), nil
}
