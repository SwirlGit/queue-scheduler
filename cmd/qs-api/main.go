package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: init config

	// TODO: init logger

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: init services

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	signal.Stop(stop)
	close(stop)
}