package storage

type Storager interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Has(key string) bool
	Count() int
}
