package shortener

import (
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/storages/filestorage"
	"github.com/ansedo/url-shortener/internal/storages/memorystorage"
	"github.com/ansedo/url-shortener/internal/storages/postgrestorage"
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

func WithPostgreStorage() Option {
	return func(s *Shortener) {
		s.Storage = helpers.Must(postgrestorage.New())
	}
}

func WithDefaultStorage() Option {
	return WithMemoryStorage()
}
