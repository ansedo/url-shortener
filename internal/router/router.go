package router

import (
	"github.com/ansedo/url-shortener/internal/handlers"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	svc := shortener.New()
	r.Post("/", handlers.EncodeURL(svc))
	r.Get("/{id}", handlers.DecodeURL(svc))
	r.Post("/api/shorten", handlers.EncodeURLFromJSON(svc))

	return r
}
