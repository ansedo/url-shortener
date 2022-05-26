package main

import (
	"github.com/ansedo/url-shortener/internal/app/config"
	"github.com/ansedo/url-shortener/internal/app/shortener"
	"github.com/ansedo/url-shortener/internal/app/storage/memory"
	"log"
	"net/http"
)

func main() {
	log.Fatal(
		http.ListenAndServe(
			config.SitePort,
			shortener.NewShortener(memory.NewStorage()),
		),
	)
}
