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
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	queryTimeout = time.Second

	softDeleteWorkersCount = 5
	softDeleteChunkSize    = 10
	softDeleteChannelSize  = 10
)

//go:embed sql
var queries embed.FS

func getQueryFromFile(filename string) (string, error) {
	query, err := queries.ReadFile(fmt.Sprintf("sql/%s", filename))
	if err != nil {
		return "", err
	}
	return string(query), nil
}

type Storage struct {
	db      *pgxpool.Pool
	queries struct {
		migrations                 string
		insertInto                 string
		selectByOriginalURL        string
		selectByShortURLID         string
		selectByUID                string
		selectMaxID                string
		softDeleteByShortURLAndUID string
	}
	chSoftDelete chan []models.ShortenLink
}

func New(ctx context.Context) (*Storage, error) {
	db, err := pgxpool.Connect(ctx, config.Get().DatabaseDSN)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		db:           db,
		chSoftDelete: make(chan []models.ShortenLink, softDeleteChannelSize),
	}

	if err = s.setQueries(ctx); err != nil {
		return nil, err
	}

	if err = s.migrate(ctx); err != nil {
		return nil, err
	}

	for i := 0; i < softDeleteWorkersCount; i++ {
		s.newSoftDeleteWorker(ctx, i+1)
	}

	s.addToShutdowner()

	return s, nil
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURLID, originalURL string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if _, err := s.db.Exec(ctx, s.queries.insertInto, helpers.GetUIDFromCtx(ctx), shortURLID, originalURL); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.SQLState()) {
				return storages.ErrRowAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (s *Storage) AddBatch(ctx context.Context, urls []models.ShortenLink) error {
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue(s.queries.insertInto, helpers.GetUIDFromCtx(ctx), url.ShortURLID, url.OriginalURL)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	br := s.db.SendBatch(ctx, batch)

	if err := br.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetBatch(ctx context.Context, urls []models.ShortenLink) ([]models.ShortenLink, error) {
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue(s.queries.selectByShortURLID, url.ShortURLID)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	br := s.db.SendBatch(ctx, batch)
	rows, err := br.Query()
	if err != nil {
		return nil, err
	}

	shortenList := make([]models.ShortenLink, len(urls))
	for rows.Next() {
		var shorten models.ShortenLink
		if err = rows.Scan(&shorten.ShortURLID, &shorten.OriginalURL); err != nil {
			return nil, err
		}
		shortenList = append(shortenList, shorten)
	}

	if err = br.Close(); err != nil {
		return nil, err
	}

	return shortenList, nil
}

func (s *Storage) AsyncSoftDeleteBatch(_ context.Context, shortenList []models.ShortenLink) {
	urlsLen := len(shortenList)
	for i := 0; i < urlsLen; i += softDeleteChunkSize {
		end := i + softDeleteChunkSize
		if end > urlsLen {
			end = urlsLen
		}
		s.chSoftDelete <- shortenList[i:end]
	}
}

func (s *Storage) SoftDeleteBatch(ctx context.Context, shortenList []models.ShortenLink) error {
	batch := &pgx.Batch{}
	for _, url := range shortenList {
		batch.Queue(s.queries.softDeleteByShortURLAndUID, url.ShortURLID, url.UID)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	br := s.db.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByShortURLID(ctx context.Context, shortURLID string) (string, error) {
	var (
		originalURL string
		isDeleted   bool
	)
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.queries.selectByShortURLID, shortURLID).Scan(&originalURL, &isDeleted); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storages.ErrShortURLIDNotExist
		}
		return "", err
	}
	if isDeleted {
		return originalURL, storages.ErrRowSoftDeleted
	}
	return originalURL, nil
}

func (s *Storage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortURLID string
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.queries.selectByOriginalURL, originalURL).Scan(&shortURLID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storages.ErrOriginalURLNotExists
		}
		return "", err
	}
	return shortURLID, nil
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenLink, error) {
	var shortenList []models.ShortenLink
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	rows, err := s.db.Query(ctx, s.queries.selectByUID, helpers.GetUIDFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var shorten models.ShortenLink
		if err = rows.Scan(&shorten.ShortURLID, &shorten.OriginalURL); err != nil {
			return nil, err
		}
		shortenList = append(shortenList, shorten)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shortenList, nil
}

func (s *Storage) GetNextID(ctx context.Context) int {
	var currentID sql.NullInt64
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.queries.selectMaxID).Scan(&currentID); err != nil {
		zap.L().Fatal(err.Error())
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
		close(s.chSoftDelete)
		s.db.Close()
		return nil
	})
}

func (s *Storage) setQueries(_ context.Context) error {
	files, err := queries.ReadDir("sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		query, err := getQueryFromFile(file.Name())
		if err != nil {
			return err
		}

		switch file.Name() {
		case "migrations.sql":
			s.queries.migrations = query
		case "insert_into.sql":
			s.queries.insertInto = query
		case "select_by_original_url.sql":
			s.queries.selectByOriginalURL = query
		case "select_by_short_url_id.sql":
			s.queries.selectByShortURLID = query
		case "select_by_uid.sql":
			s.queries.selectByUID = query
		case "select_max_id.sql":
			s.queries.selectMaxID = query
		case "soft_delete_by_short_url_and_uid.sql":
			s.queries.softDeleteByShortURLAndUID = query
		}
	}
	return nil
}

func (s *Storage) migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	_, err := s.db.Exec(ctx, s.queries.migrations)
	return err
}

func (s *Storage) newSoftDeleteWorker(ctx context.Context, id int) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				zap.L().Warn("soft delete worker context is canceled",
					zap.String("soft_delete_worker_id", strconv.Itoa(id)),
				)
				return
			case urls := <-s.chSoftDelete:
				s.doSoftDeleteWork(ctx, id, urls)
			}
		}
	}()
}

func (s *Storage) doSoftDeleteWork(ctx context.Context, id int, urls []models.ShortenLink) {
	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	if err := s.SoftDeleteBatch(ctx, urls); err != nil {
		zap.L().Warn(err.Error(),
			zap.String("method", "newSoftDeleteWorker"),
			zap.String("soft_delete_worker_id", strconv.Itoa(id)),
		)
	}
}
