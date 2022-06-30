package main

import (
	"github.com/ansedo/url-shortener/internal/server"
	"github.com/ansedo/url-shortener/internal/services/shutdowner"
)

func main() {
	server.Run()
	<-shutdowner.Get().ChShutdowned
}
