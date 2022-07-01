package router

import (
	"context"
	"github.com/ansedo/url-shortener/internal/handlers"
	"github.com/ansedo/url-shortener/internal/middlewares"
	"github.com/ansedo/url-shortener/internal/services/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(ctx context.Context) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.Decompress)
	r.Use(middlewares.Cookie)

	svc := shortener.New(ctx)

	r.Post("/", handlers.ShortenURL(svc))
	r.Get("/{id}", handlers.GetOriginalURL(svc))
	r.Get("/ping", handlers.PingDB(svc))

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.APIShortenURL(svc))
		r.Get("/user/urls", handlers.APIGetURLsByUID(svc))
	})

	return r
}
