package logger

import (
	"context"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"go.uber.org/zap"
	"sync"
)

var once sync.Once

func New() {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			zap.L().Fatal(err.Error())
		}
		addToShutdowner(logger)
		zap.ReplaceGlobals(logger)
	})
}

func addToShutdowner(logger *zap.Logger) {
	shutdowner.Get().AddCloser(func(ctx context.Context) error {
		if err := logger.Sync(); err != nil {
			return err
		}
		return nil
	})
}
