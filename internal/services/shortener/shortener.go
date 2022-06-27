package shortener

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/storages"
	"github.com/speps/go-hashids/v2"
)

const (
	hashSalt   = "The Little Man Who Wasn't There"
	hashLength = 8
)

type Shortener struct {
	Storage storages.Storager
}

func New(opts ...Option) *Shortener {
	s := &Shortener{}

	for _, opt := range opts {
		opt(s)
	}

	if s.Storage == nil && config.Get().FileStoragePath != "" {
		WithFileStorage()(s)
	}

	if s.Storage == nil {
		WithDefaultStorage()(s)
	}

	return s
}

func (s *Shortener) GenerateID() (string, error) {
	d := hashids.NewData()
	d.Salt = hashSalt
	d.MinLength = hashLength
	h, err := hashids.NewWithData(d)
	if err != nil {
		return "", err
	}

	id, err := h.Encode([]int{s.Storage.NextID()})
	if err != nil {
		return "", err
	}

	return id, nil
}
