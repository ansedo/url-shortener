package shortener

import (
	"github.com/ansedo/url-shortener/internal/storage/memorystorage"
)

type Option func(s *Shortener)

func WithMemoryStorage() Option {
	return func(s *Shortener) {
		s.Storage = memorystorage.New()
	}
}

func WithDefaultStorage() Option {
	return WithMemoryStorage()
}
