package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"sync"
)

var once sync.Once

type config struct {
	BaseURL                string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress          string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	RequestNotAllowedError string `env:"REQUEST_NOT_ALLOWED_ERROR" envDefault:"this request is not allowed"`
}

var instance *config

func New(opts ...Option) *config {
	if instance == nil {
		once.Do(
			func() {
				instance = &config{}
				err := env.Parse(instance)
				if err != nil {
					log.Fatal(err)
				}

				for _, opt := range opts {
					opt(instance)
				}
			})
	}
	return instance
}
