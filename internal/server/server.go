package server

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/router"
	"log"
	"net/http"
)

func Run() *http.Server {
	srv := &http.Server{
		Addr:    config.NewConfig().SitePort,
		Handler: router.NewRouter(),
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}
