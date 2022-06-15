package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"sync"
)

var once sync.Once

type config struct {
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var instance *config

func New(opts ...Option) *config {
	if instance == nil {
		once.Do(
			func() {
				var cfg config
				err := env.Parse(&cfg)
				if err != nil {
					log.Fatal(err)
				}

				flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, `server address to listen on`)
				flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, `basic URL of resulting shortened URL`)
				flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, `location to store data in`)
				flag.Parse()

				for _, opt := range opts {
					opt(&cfg)
				}

				instance = &cfg
			})
	}
	return instance
}
