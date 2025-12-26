package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	log.Println("authforge: Server starting")

	<-signalChan
	log.Println("authforge: Server interuptted")

	cancel()
	log.Println("authforge: Server shutting down")
}
