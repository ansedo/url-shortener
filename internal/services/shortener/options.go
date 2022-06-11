package shortener

import (
	"github.com/ansedo/url-shortener/internal/storages/filestorage"
	"github.com/ansedo/url-shortener/internal/storages/memorystorage"
)

type Option func(s *Shortener)

func WithMemoryStorage() Option {
	return func(s *Shortener) {
		s.Storage = memorystorage.New()
	}
}

func WithFileStorage() Option {
	return func(s *Shortener) {
		s.Storage = filestorage.New()
	}
}

func WithDefaultStorage() Option {
	return WithMemoryStorage()
}
