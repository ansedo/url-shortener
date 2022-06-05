package shortener

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/storage/memory"
)

type Option func(s *Shortener)

func WithConfig(cfg *config.Config) Option {
	return func(s *Shortener) {
		s.Config = cfg
	}
}

func WithDefaultConfig() Option {
	return func(s *Shortener) {
		s.Config = config.NewConfig()
	}
}

func WithMemoryStorage() Option {
	return func(s *Shortener) {
		s.Storage = memory.NewStorage()
	}
}

func WithDefaultStorage() Option {
	return WithMemoryStorage()
}
