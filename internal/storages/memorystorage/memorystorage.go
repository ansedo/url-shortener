package memorystorage

import (
	"context"
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

func New(_ context.Context) *Storage {
	return &Storage{
		repo: make(map[string]row),
	}
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURLID, originalURL string) error {
	if s.IsDuplicate(ctx, shortURLID, originalURL) {
		return storages.ErrRowAlreadyExists
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo[shortURLID] = row{
		UID:         helpers.GetUIDFromCtx(ctx),
		OriginalURL: originalURL,
	}
	return nil
}

func (s *Storage) AddBatch(ctx context.Context, urls []models.ShortenList) error {
	for _, url := range urls {
		if err := s.Add(ctx, url.ShortURLID, url.OriginalURL); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetByShortURLID(_ context.Context, shortURLID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.repo[shortURLID]
	if !ok {
		return "", storages.ErrShortURLIDNotExist
	}
	return row.OriginalURL, nil
}

func (s *Storage) GetByOriginalURL(_ context.Context, originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, row := range s.repo {
		if row.OriginalURL == originalURL {
			return originalURL, nil
		}
	}
	return "", storages.ErrOriginalURLNotExists
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenList, error) {
	var shortenList []models.ShortenList
	uid := helpers.GetUIDFromCtx(ctx)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for shortURLID, row := range s.repo {
		if row.UID == uid {
			shortenList = append(
				shortenList,
				models.ShortenList{
					ShortURLID:  shortURLID,
					OriginalURL: row.OriginalURL,
				},
			)
		}
	}
	return shortenList, nil
}

func (s *Storage) IsDuplicate(ctx context.Context, shortURLID, originalURL string) bool {
	if _, err := s.GetByShortURLID(ctx, shortURLID); err == nil {
		return true
	}
	if _, err := s.GetByOriginalURL(ctx, originalURL); err == nil {
		return true
	}
	return false
}

func (s *Storage) GetNextID(_ context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.repo)
}

func (s *Storage) Ping(_ context.Context) error {
	return nil
}
