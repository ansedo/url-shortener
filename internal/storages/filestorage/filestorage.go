package filestorage

import (
	"context"
	"errors"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/helpers"
	"github.com/ansedo/url-shortener/internal/models"
	"github.com/ansedo/url-shortener/internal/storages"
	"io"
	"log"
)

type Storage struct {
	fileName string
}

type Record struct {
	UID         string `json:"uid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func New() *Storage {
	return &Storage{
		fileName: config.Get().FileStoragePath,
	}
}

var _ storages.Storager = (*Storage)(nil)

func (s *Storage) Add(ctx context.Context, shortURL, originalURL string) error {
	if s.IsShortURLExist(ctx, shortURL) {
		return storages.ErrKeyAlreadyExists
	}
	producer := helpers.Must(NewProducer(s.fileName))
	defer producer.Close()
	err := producer.WriteRecord(&Record{
		UID:         helpers.GetUIDFromCtx(ctx),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetByShortURL(_ context.Context, shortURL string) (string, error) {
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
		if record.ShortURL == shortURL {
			return record.OriginalURL, nil
		}
	}
	return "", storages.ErrKeyNotExist
}

func (s *Storage) GetByUID(ctx context.Context) ([]models.ShortenListResponse, error) {
	entities := make([]models.ShortenListResponse, 0)
	uid := helpers.GetUIDFromCtx(ctx)
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record.UID == uid {
			entities = append(
				entities,
				models.ShortenListResponse{
					ShortURL:    record.ShortURL,
					OriginalURL: record.OriginalURL,
				})
		}
	}
	return entities, nil
}

func (s *Storage) IsShortURLExist(_ context.Context, shortURL string) bool {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record.ShortURL == shortURL {
			return true
		}
	}
	return false
}

func (s *Storage) NextID(_ context.Context) int {
	consumer := helpers.Must(NewConsumer(s.fileName))
	defer consumer.Close()
	var nextID int
	for {
		if _, err := consumer.ReadRecord(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}
		nextID++
	}
	return nextID + 1
}
