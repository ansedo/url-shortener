package server

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/router"
	"log"
	"net/http"
)

func Run() *http.Server {
	srv := &http.Server{
		Addr:    config.New().SitePort,
		Handler: router.New(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return srv
}
