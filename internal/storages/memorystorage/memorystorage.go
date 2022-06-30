package memorystorage

import (
	"context"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/storages"
	"sync"
)

type row struct {
	UID         string
	OriginalURL string
}

type Storage struct {
	mu   sync.RWMutex
	repo map[string]row
}

func New() *Storage {
	return &Storage{
		repo: make(map[string]row),
	}
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURL, originalURL string) error {
	if s.IsShortURLExist(ctx, shortURL) {
		return storages.ErrKeyAlreadyExists
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo[shortURL] = row{
		UID:         helpers.GetUIDFromCtx(ctx),
		OriginalURL: originalURL,
	}
	return nil
}

func (s *Storage) GetByShortURL(ctx context.Context, shortURL string) (string, error) {
	if !s.IsShortURLExist(ctx, shortURL) {
		return "", storages.ErrKeyNotExist
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repo[shortURL].OriginalURL, nil
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenListResponse, error) {
	entities := make([]models.ShortenListResponse, 0)
	uid := helpers.GetUIDFromCtx(ctx)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for shortURL, row := range s.repo {
		if row.UID == uid {
			entities = append(
				entities,
				models.ShortenListResponse{
					ShortURL:    config.Get().BaseURL + "/" + shortURL,
					OriginalURL: row.OriginalURL,
				})
		}
	}
	return entities, nil
}

func (s *Storage) IsShortURLExist(_ context.Context, shortURL string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.repo[shortURL]
	return ok
}

func (s *Storage) NextID(_ context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.repo)
}

func (s *Storage) Ping(_ context.Context) error {
	return nil
}
