package shortener

import (
	"context"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/storages/filestorage"
	"github.com/ansedo/url-shortener/internal/storages/memorystorage"
	"github.com/ansedo/url-shortener/internal/storages/postgrestorage"
)

type Option func(s *Shortener)

func WithMemoryStorage(ctx context.Context) Option {
	return func(s *Shortener) {
		s.Storage = memorystorage.New(ctx)
	}
}

func WithFileStorage(ctx context.Context) Option {
	return func(s *Shortener) {
		s.Storage = filestorage.New(ctx)
	}
}

func WithPostgreStorage(ctx context.Context) Option {
	return func(s *Shortener) {
		s.Storage = helpers.Must(postgrestorage.New(ctx))
	}
}

func WithDefaultStorage(ctx context.Context) Option {
	return WithMemoryStorage(ctx)
}
