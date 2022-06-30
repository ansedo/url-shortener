package server

import (
	"context"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/router"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"log"
	"net/http"
)

type Server struct {
	http http.Server
}

func New() *Server {
	return &Server{
		http: http.Server{
			Addr:    config.Get().ServerAddress,
			Handler: router.New(),
		},
	}
}

func Run() {
	srv := New()
	go srv.ListenAndServer()
	srv.addToShutdowner()
}

func (s *Server) ListenAndServer() {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *Server) addToShutdowner() {
	shutdowner.Get().AddCloser(func(ctx context.Context) error {
		if err := s.http.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
}
