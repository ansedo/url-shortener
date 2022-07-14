package main

import (
	"context"
	"github.com/ansedo/url-shortener/internal/logger"
	"github.com/ansedo/url-shortener/internal/server"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
	"sync"
)

var once sync.Once

func main() {
	logger.New()
	server.Run(context.Background())
	<-shutdowner.Get().ChShutdowned
}
