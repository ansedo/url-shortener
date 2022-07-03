package storages

import "errors"

var (
	ErrShortURLIDAlreadyExists = errors.New("this short url id already exists")
	ErrShortURLIDNotExist      = errors.New("this short url id does not exist")
)
