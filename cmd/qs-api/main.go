package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwirlGit/queue-scheduler/internal/qs-api/api/v1/schedule"
	"github.com/SwirlGit/queue-scheduler/pkg/fasthttp"
)

func main() {
	// TODO: init config

	// TODO: init logger

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: init services
	scheduleService := schedule.NewService()
	scheduleHandler := schedule.NewHandler(scheduleService)

	server := fasthttp.NewServer([]fasthttp.RouteProvider{scheduleHandler})
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
