package shortener

import (
	"context"
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
	BaseURL string
	NextID  int
}

func New(ctx context.Context, opts ...Option) *Shortener {
	s := &Shortener{
		BaseURL: config.Get().BaseURL,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.Storage == nil && config.Get().DatabaseDSN != "" {
		WithPostgreStorage(ctx)(s)
	}

	if s.Storage == nil && config.Get().FileStoragePath != "" {
		WithFileStorage(ctx)(s)
	}

	if s.Storage == nil {
		WithDefaultStorage(ctx)(s)
	}

	s.NextID = s.Storage.NextID(ctx)

	return s
}

func (s *Shortener) GenerateID(_ context.Context) (string, error) {
	d := hashids.NewData()
	d.Salt = hashSalt
	d.MinLength = hashLength
	h, err := hashids.NewWithData(d)
	if err != nil {
		return "", err
	}

	id, err := h.Encode([]int{s.NextID})
	if err != nil {
		return "", err
	}
	s.NextID++

	return id, nil
}
