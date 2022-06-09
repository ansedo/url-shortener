package config

import (
	"sync"
)

var once sync.Once

type config struct {
	SiteScheme  string
	SiteHost    string
	SitePort    string
	SiteAddress string

	RequestNotAllowedError string
}

var instance *config

func New(opts ...Option) *config {
	if instance == nil {
		once.Do(
			func() {
				instance = &config{
					SiteScheme: "http://",
					SiteHost:   "localhost",
					SitePort:   ":8080",

					RequestNotAllowedError: "this request is not allowed",
				}

				for _, opt := range opts {
					opt(instance)
				}

				instance.SiteAddress = instance.SiteScheme + instance.SiteHost + instance.SitePort
			})
	}
	return instance
}
