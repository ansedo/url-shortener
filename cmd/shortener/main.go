package main

import (
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/handlers"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	shortenerService := shortener.NewShortener()
	router.Post("/", handlers.EncodeURL(shortenerService))
	router.Get("/{id}", handlers.DecodeURL(shortenerService))

	log.Fatal(http.ListenAndServe(config.NewConfig().SitePort, router))
}
