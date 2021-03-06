package filestorage

import (
	"context"
	"errors"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/storages"
	"go.uber.org/zap"
	"io"
)

type Storage struct {
	fileName string
}

type Record struct {
	UID         string `json:"uid"`
	ShortURLID  string `json:"short_url_id"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"is_deleted"`
}

func New(_ context.Context) *Storage {
	return &Storage{
		fileName: config.Get().FileStoragePath,
	}
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURLID, originalURL string) error {
	if s.IsDuplicate(ctx, shortURLID, originalURL) {
		return storages.ErrRowAlreadyExists
	}
	producer := helpers.Must(NewProducer(s.fileName))
	defer producer.Close()
	err := producer.WriteRecord(&Record{
		UID:         helpers.GetUIDFromCtx(ctx),
		ShortURLID:  shortURLID,
		OriginalURL: originalURL,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddBatch(ctx context.Context, urls []models.ShortenLink) error {
	for _, url := range urls {
		if err := s.Add(ctx, url.ShortURLID, url.OriginalURL); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) AsyncSoftDeleteBatch(ctx context.Context, urls []models.ShortenLink) {
	s.SoftDeleteBatch(ctx, urls)
}

func (s *Storage) SoftDeleteBatch(ctx context.Context, urls []models.ShortenLink) {
	repo := make(map[string]Record)
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			zap.L().Warn(err.Error(),
				zap.String("method", "SoftDeleteBatch"),
			)
			break
		}
		repo[record.ShortURLID] = *record
	}

	for _, url := range urls {
		record, ok := repo[url.ShortURLID]
		if !ok {
			continue
		}
		if repo[url.ShortURLID].UID == url.UID {
			record.IsDeleted = true
			repo[url.ShortURLID] = record
		}
	}

	for _, record := range repo {
		if err := s.Add(ctx, record.ShortURLID, record.OriginalURL); err != nil {
			zap.L().Warn(err.Error(),
				zap.String("method", "SoftDeleteBatch"),
			)
		}
	}
}

func (s *Storage) GetByShortURLID(_ context.Context, shortURLID string) (string, error) {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if record.ShortURLID == shortURLID {
			if record.IsDeleted {
				return record.OriginalURL, storages.ErrRowSoftDeleted
			}
			return record.OriginalURL, nil
		}
	}
	return "", storages.ErrShortURLIDNotExist
}

func (s *Storage) GetByOriginalURL(_ context.Context, originalURL string) (string, error) {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if record.ShortURLID == originalURL {
			return record.ShortURLID, nil
		}
	}
	return "", storages.ErrOriginalURLNotExists
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenLink, error) {
	entities := make([]models.ShortenLink, 0)
	uid := helpers.GetUIDFromCtx(ctx)
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if record.UID == uid {
			entities = append(
				entities,
				models.ShortenLink{
					ShortURLID:  record.ShortURLID,
					OriginalURL: record.OriginalURL,
				})
		}
	}
	return entities, nil
}

func (s *Storage) IsDuplicate(_ context.Context, shortURLID, originalURL string) bool {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		if record.ShortURLID == shortURLID || record.OriginalURL == originalURL {
			return false
		}
	}
	return false
}

func (s *Storage) GetNextID(_ context.Context) int {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	var nextID int
	for {
		if _, err := consumer.ReadRecord(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			zap.L().Fatal(err.Error())
		}
		nextID++
	}
	return nextID + 1
}

func (s *Storage) Ping(_ context.Context) error {
	return nil
}
