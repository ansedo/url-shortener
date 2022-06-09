package shortener

import (
	"github.com/ansedo/url-shortener/internal/storage/memorystorage"
)

type Option func(s *Shortener)

func WithMemoryStorage() Option {
	return func(s *Shortener) {
		s.Storage = memorystorage.NewStorage()
	}
}

func WithDefaultStorage() Option {
	return WithMemoryStorage()
}
