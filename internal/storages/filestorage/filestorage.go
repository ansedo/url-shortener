package filestorage

import (
	"errors"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/storages"
	"io"
	"log"
)

type Storage struct {
	fileName string
}

type Record struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func New() *Storage {
	return &Storage{
		fileName: config.Get("FileStoragePath"),
	}
}

func (s *Storage) Get(key string) (string, error) {
	consumer, err := NewConsumer(s.fileName)
	if err != nil {
		return "", err
	}
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if record.ID == key {
			return record.URL, nil
		}
	}
	return "", storages.ErrKeyNotExist
}

func (s *Storage) Set(key, value string) error {
	if s.Has(key) {
		return storages.ErrKeyAlreadyExists
	}
	producer, err := NewProducer(s.fileName)
	if err != nil {
		return err
	}
	defer producer.Close()
	err = producer.WriteRecord(&Record{ID: key, URL: value})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Has(key string) bool {
	consumer, err := NewConsumer(s.fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	for {
		record, err := consumer.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record.ID == key {
			return true
		}
	}
	return false
}

func (s *Storage) NextID() int {
	consumer, err := NewConsumer(s.fileName)
	if err != nil {
		log.Fatal(err)
	}
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
	log.Println(nextID)
	return nextID + 1
}
