package storages

import (
	"context"
	"github.com/ansedo/url-shortener/internal/models"
)

// Storager is the common interface implemented by all storages.
type Storager interface {
	Add(ctx context.Context, shortURLID, originalURL string) error
	AddBatch(ctx context.Context, urls []models.ShortenList) error
	GetByShortURLID(ctx context.Context, shortURLID string) (string, error)
	GetByUID(ctx context.Context) ([]models.ShortenList, error)
	IsShortURLIDExist(ctx context.Context, shortURLID string) bool
	NextID(ctx context.Context) int
	Ping(ctx context.Context) error
}
