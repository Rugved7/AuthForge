package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rugved7/authforge/internal/config"
	apphttp "github.com/Rugved7/authforge/internal/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	// create router
	router := apphttp.NewRouter()

	// create http server
	server := apphttp.NewServer(cfg.ServicePort, router)

	// start server
	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("authforge: Server starting")

	<-signalChan
	log.Println("authforge: Server interuptted")

	cancel()
	log.Println("authforge: Server shutting down")
}
