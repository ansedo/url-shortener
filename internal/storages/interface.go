package storages

import (
	"context"
	"github.com/ansedo/url-shortener/internal/models"
)

// Storager is the common interface implemented by all storages.
type Storager interface {
	Add(ctx context.Context, shortURL, originalURL string) error
	GetByShortURL(ctx context.Context, shortURL string) (string, error)
	GetByUID(ctx context.Context) ([]models.ShortenListResponse, error)
	IsShortURLExist(ctx context.Context, shortURL string) bool
	NextID(ctx context.Context) int
}
