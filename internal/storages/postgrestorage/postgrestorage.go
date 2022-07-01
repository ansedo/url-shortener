package postgrestorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"github.com/ansedo/url-shortener/internal/storages"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

const queryTimeout = time.Second

//go:embed sql
var queries embed.FS

func getSQL(name string) string {
	return string(helpers.Must(queries.ReadFile(fmt.Sprintf("sql/%s.sql", name))))
}

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context) (*Storage, error) {
	db, err := pgx.Connect(ctx, config.Get().DatabaseDSN)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		db: db,
	}

	err = s.migrate(ctx)
	if err != nil {
		return nil, err
	}

	s.addToShutdowner()

	return s, nil
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURL, originalURL string) error {
	if s.IsShortURLExist(ctx, shortURL) {
		return storages.ErrShortURLAlreadyExists
	}
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	_, err := s.db.Exec(ctx, getSQL("insert_into_urls"), helpers.GetUIDFromCtx(ctx), shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetByShortURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	err := s.db.QueryRow(ctx, getSQL("select_by_short_url"), shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storages.ErrShortURLNotExist
		}
		return "", err
	}
	return originalURL, nil
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenListResponse, error) {
	var shortenList []models.ShortenListResponse
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	rows, err := s.db.Query(ctx, getSQL("select_by_uid"), helpers.GetUIDFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var entity models.ShortenListResponse
		err = rows.Scan(&entity.ShortURL, &entity.OriginalURL)
		if err != nil {
			return nil, err
		}
		shortenList = append(shortenList, entity)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return shortenList, nil
}

func (s *Storage) IsShortURLExist(ctx context.Context, shortURL string) bool {
	var isExist bool
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	err := s.db.QueryRow(ctx, getSQL("select_exists_by_short_url"), shortURL).Scan(&isExist)
	if err != nil {
		log.Fatal(err)
	}
	return isExist
}

func (s *Storage) NextID(ctx context.Context) int {
	var currentID sql.NullInt64
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	err := s.db.QueryRow(ctx, getSQL("select_max_id")).Scan(&currentID)
	if err != nil {
		log.Fatal(err)
	}
	if !currentID.Valid {
		return 0
	}
	return int(currentID.Int64) + 1
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

func (s *Storage) migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	_, err := s.db.Exec(ctx, getSQL("migrate"))
	return err
}
