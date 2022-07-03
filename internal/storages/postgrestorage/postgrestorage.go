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
	"go.uber.org/zap"
	"time"
)

const queryTimeout = time.Second

//go:embed sql/migrations.sql
var migrations string

//go:embed sql/insert* sql/select*
var queries embed.FS

func getQueryFromFile(filename string) (string, error) {
	query, err := queries.ReadFile(fmt.Sprintf("sql/%s", filename))
	if err != nil {
		return "", err
	}
	return string(query), nil
}

type Storage struct {
	db    *pgx.Conn
	stmts struct {
		insertInto          *pgconn.StatementDescription
		selectByOriginalURL *pgconn.StatementDescription
		selectByShortURLID  *pgconn.StatementDescription
		selectByUID         *pgconn.StatementDescription
		selectMaxID         *pgconn.StatementDescription
	}
}

func New(ctx context.Context) (*Storage, error) {
	db, err := pgx.Connect(ctx, config.Get().DatabaseDSN)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		db: db,
	}

	if err = s.migrate(ctx); err != nil {
		return nil, err
	}

	if err = s.setStatements(ctx); err != nil {
		return nil, err
	}

	s.addToShutdowner()

	return s, nil
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURLID, originalURL string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if _, err := s.db.Exec(ctx, s.stmts.insertInto.Name, helpers.GetUIDFromCtx(ctx), shortURLID, originalURL); err != nil {
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

func (s *Storage) AddBatch(ctx context.Context, urls []models.ShortenList) error {
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue(s.stmts.insertInto.Name, helpers.GetUIDFromCtx(ctx), url.ShortURLID, url.OriginalURL)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	br := s.db.SendBatch(ctx, batch)

	if err := br.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetBatch(ctx context.Context, urls []models.ShortenList) ([]models.ShortenList, error) {
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue(s.stmts.selectByShortURLID.Name, url.ShortURLID)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	br := s.db.SendBatch(ctx, batch)
	rows, err := br.Query()
	if err != nil {
		return nil, err
	}

	shortenList := make([]models.ShortenList, len(urls))
	for rows.Next() {
		var shorten models.ShortenList
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

func (s *Storage) GetByShortURLID(ctx context.Context, shortURLID string) (string, error) {
	var originalURL string
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.stmts.selectByShortURLID.Name, shortURLID).Scan(&originalURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storages.ErrShortURLIDNotExist
		}
		return "", err
	}
	return originalURL, nil
}

func (s *Storage) GetByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortURLID string
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.stmts.selectByOriginalURL.Name, originalURL).Scan(&shortURLID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storages.ErrOriginalURLNotExists
		}
		return "", err
	}
	return shortURLID, nil
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenList, error) {
	var shortenList []models.ShortenList
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	rows, err := s.db.Query(ctx, s.stmts.selectByUID.Name, helpers.GetUIDFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var shorten models.ShortenList
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

func (s *Storage) NextID(ctx context.Context) int {
	var currentID sql.NullInt64
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.QueryRow(ctx, s.stmts.selectMaxID.Name).Scan(&currentID); err != nil {
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
		if err := s.db.Close(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) setStatements(ctx context.Context) error {
	files, err := queries.ReadDir("sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		query, err := getQueryFromFile(file.Name())
		if err != nil {
			return err
		}

		stmt, err := s.db.Prepare(ctx, file.Name(), query)
		if err != nil {
			return err
		}

		switch file.Name() {
		case "insert_into.sql":
			s.stmts.insertInto = stmt
		case "select_by_original_url.sql":
			s.stmts.selectByOriginalURL = stmt
		case "select_by_short_url_id.sql":
			s.stmts.selectByShortURLID = stmt
		case "select_by_uid.sql":
			s.stmts.selectByUID = stmt
		case "select_max_id.sql":
			s.stmts.selectMaxID = stmt
		}
	}
	return nil
}

func (s *Storage) migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	_, err := s.db.Exec(ctx, migrations)
	return err
}
