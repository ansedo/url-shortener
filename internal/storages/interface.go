package storages

import (
	"context"
	"github.com/ansedo/url-shortener/internal/models"
)

// Storager is the common interface implemented by all storages.
type Storager interface {
	StorageReader
	StorageWriter
	AsyncSoftDeleteBatch(ctx context.Context, urls []models.ShortenLink)
	Ping(ctx context.Context) error
}

type StorageReader interface {
	GetByShortURLID(ctx context.Context, shortURLID string) (string, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (string, error)
	GetByUID(ctx context.Context) ([]models.ShortenLink, error)
	GetNextID(ctx context.Context) int
}

type StorageWriter interface {
	Add(ctx context.Context, shortURLID, originalURL string) error
	AddBatch(ctx context.Context, urls []models.ShortenLink) error
}
