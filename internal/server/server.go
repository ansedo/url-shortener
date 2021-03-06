package server

import (
	"context"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/router"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	http http.Server
}

func New(ctx context.Context) *Server {
	return &Server{
		http: http.Server{
			Addr:    config.Get().ServerAddress,
			Handler: router.New(ctx),
		},
	}
}

func Run(ctx context.Context) {
	srv := New(ctx)
	go srv.ListenAndServer()
	srv.addToShutdowner()
}

func (s *Server) ListenAndServer() {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.L().Fatal(err.Error())
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
