package postgrestorage

import (
	"context"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"github.com/ansedo/url-shortener/internal/storages"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	db *pgx.Conn
}

func New() (*Storage, error) {
	db, err := pgx.Connect(context.Background(), config.Get().DatabaseDSN)
	if err != nil {
		return nil, err
	}
	s := &Storage{
		db: db,
	}
	s.addToShutdowner()
	return s, nil
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(_ context.Context, _, _ string) error {
	return nil
}

func (s *Storage) GetByShortURL(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (s *Storage) GetByUID(_ context.Context) ([]models.ShortenListResponse, error) {
	return nil, nil
}

func (s *Storage) IsShortURLExist(_ context.Context, _ string) bool {
	return false
}

func (s *Storage) NextID(_ context.Context) int {
	return 0
}

func (s *Storage) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Storage) addToShutdowner() {
	shutdowner.Get().AddCloser(func(ctx context.Context) error {
		err := s.db.Close(context.Background())
		if err != nil {
			return err
		}
		return nil
	})
}
