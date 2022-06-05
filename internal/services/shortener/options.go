package shortener

import "github.com/ansedo/url-shortener/internal/storage/memory"

type Option func(s *Shortener)

func WithMemoryStorage() Option {
	return func(s *Shortener) {
		s.Storage = memory.NewStorage()
	}
}
