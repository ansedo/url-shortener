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
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	svc := shortener.NewShortener()
	r.Post("/", handlers.EncodeURL(svc))
	r.Get("/{id}", handlers.DecodeURL(svc))

	log.Fatal(http.ListenAndServe(config.NewConfig().SitePort, r))
}
