package main

import (
	"context"
	"fmt"
	"jayant-001/api-gateway/internal/config"
	"jayant-001/api-gateway/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	r, err := router.NewRouter(config)
	if err != nil {
		fmt.Printf("Error creating router: %v\n", err)
		return
	}

	srv := http.Server{
		Addr : config.Server.PORT,
		Handler: r,
		ReadTimeout: config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	go func() {
		log.Printf("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<- quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	
	log.Println("server gracefully stopped")
}
