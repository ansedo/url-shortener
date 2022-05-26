package storage

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Has(key string) bool
	Count() int
}
