package main

import (
	"context"
	"github.com/ansedo/url-shortener/internal/server"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
)

func main() {
	server.Run(context.Background())
	<-shutdowner.Get().ChShutdowned
}
