package main

import (
	"github.com/ansedo/url-shortener/internal/app/config"
	"github.com/ansedo/url-shortener/internal/app/shortener"
	"log"
	"net/http"
)

func main() {
	log.Fatal(
		http.ListenAndServe(
			config.SitePort,
			shortener.NewShortener(),
		),
	)
}
