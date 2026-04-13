package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sakib-maho/Golang-Beggo-RestAPI-Swagger/internal/api"
	"github.com/sakib-maho/Golang-Beggo-RestAPI-Swagger/internal/config"
	"github.com/sakib-maho/Golang-Beggo-RestAPI-Swagger/internal/store"
)

func main() {
	cfg := config.Load()
	taskStore := store.NewMemoryTaskStore()
	handler := api.NewHandler(taskStore)

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           api.NewRouter(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
