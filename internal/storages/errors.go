package storages

import "errors"

var (
	ErrKeyAlreadyExists = errors.New("this key already exists")
	ErrKeyNotExist      = errors.New("this key does not exist")
)
