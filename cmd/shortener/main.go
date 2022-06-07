package main

import (
	"context"
	"github.com/ansedo/url-shortener/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	waitShutdownServerMaxSeconds = 5
)

func main() {
	srv := server.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), waitShutdownServerMaxSeconds*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
