package filestorage

import (
	"errors"
	"github.com/ansedo/url-shortener/internal/config"
	"io"
	"log"
	"strconv"
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
		fileName: config.New().FileStoragePath,
	}
}

func (s *Storage) Get(key string) (string, error) {
	consumer, err := NewConsumer(s.fileName)
	defer consumer.Close()
	if err != nil {
		return "", err
	}
	for {
		record, err := consumer.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if record.ID == key {
			return record.URL, nil
		}
	}
	return "", errors.New("this key does not exist")
}

func (s *Storage) Set(key, value string) error {
	if s.Has(key) {
		return errors.New("this key already exists")
	}
	producer, err := NewProducer(s.fileName)
	defer producer.Close()
	if err != nil {
		return err
	}
	err = producer.WriteRecord(&Record{ID: key, URL: value})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Has(key string) bool {
	consumer, err := NewConsumer(s.fileName)
	defer consumer.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		record, err := consumer.ReadRecord()
		if err == io.EOF {
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
	defer consumer.Close()
	if err != nil {
		log.Fatal(err)
	}

	record, err := consumer.ReadLastRecord()
	if err != nil {
		log.Fatal(err)
	}

	if record == nil {
		return 0
	}

	nextID, err := strconv.Atoi(record.ID)
	if err != nil {
		log.Fatal(err)
	}

	return nextID + 1
}
