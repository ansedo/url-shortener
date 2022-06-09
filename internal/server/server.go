package server

import (
	"context"
	"errors"
	"github.com/ansedo/url-shortener/internal/config"
	"github.com/ansedo/url-shortener/internal/router"
	"log"
	"net/http"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second

	gracefulShutdownMessage = "server gracefully finished"
)

func Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    config.New().SitePort,
		Handler: router.NewRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return errors.New(gracefulShutdownMessage)
}
