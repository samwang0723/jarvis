package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"samwang0723/jarvis/handlers"
	"samwang0723/jarvis/services"
	"syscall"
	"time"
)

func main() {
	handler := handlers.New(services.NewService())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	handler.BatchingDownload(context.Background(), -1, 5000)

	<-done
	log.Println("server stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	log.Println("server exited properly")
}
