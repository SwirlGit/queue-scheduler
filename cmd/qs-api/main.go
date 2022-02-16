package main

import (
	"context"
	"github.com/queue-scheduler/pkg/fasthttp"
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

	server := fasthttp.NewServer(nil)
	go func() {
		if err := server.Listen(":9000"); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	signal.Stop(stop)
	close(stop)

	if err := server.Shutdown(); err != nil {
		panic(err)
	}
}