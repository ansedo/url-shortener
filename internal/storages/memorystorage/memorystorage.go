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
	sync.RWMutex
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
	s.Lock()
	defer s.Unlock()
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
	s.RLock()
	defer s.RUnlock()
	return s.repo[shortURL].OriginalURL, nil
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenListResponse, error) {
	entities := make([]models.ShortenListResponse, 0)
	uid := helpers.GetUIDFromCtx(ctx)
	s.RLock()
	defer s.RUnlock()
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
	s.RLock()
	defer s.RUnlock()
	_, ok := s.repo[shortURL]
	return ok
}

func (s *Storage) NextID(_ context.Context) int {
	s.RLock()
	defer s.RUnlock()
	return len(s.repo)
}
