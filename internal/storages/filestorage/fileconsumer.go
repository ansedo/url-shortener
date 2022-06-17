package filestorage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadRecord() (*Record, error) {
	var record Record
	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *consumer) ReadLastRecord() (*Record, error) {
	stat, err := c.file.Stat()
	if err != nil {
		return nil, err
	}
	filesize := stat.Size()
	if filesize == 0 {
		return nil, nil
	}

	var line string
	var cursor int64
	for {
		cursor -= 1
		_, err = c.file.Seek(cursor, io.SeekEnd)
		if err != nil {
			return nil, err
		}

		char := make([]byte, 1)
		_, err = c.file.Read(char)
		if err != nil {
			return nil, err
		}

		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}

		line = fmt.Sprintf("%s%s", char, line)

		if cursor == -filesize {
			break
		}
	}

	if line == "" {
		return nil, nil
	}

	var record Record
	if err = json.Unmarshal([]byte(line), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
