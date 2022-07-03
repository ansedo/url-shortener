package storages

import "errors"

var (
	ErrRowAlreadyExists     = errors.New("this row already exists")
	ErrShortURLIDNotExist   = errors.New("this short url id does not exist")
	ErrOriginalURLNotExists = errors.New("this original url does not exist")
)
