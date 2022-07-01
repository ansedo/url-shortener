package storages

import "errors"

var (
	ErrShortURLAlreadyExists = errors.New("this short url already exists")
	ErrShortURLNotExist      = errors.New("this short url does not exist")
)
