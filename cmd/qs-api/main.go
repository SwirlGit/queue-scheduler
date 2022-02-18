package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwirlGit/queue-scheduler/cmd/qs-api/config"
	pkgschedule "github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/SwirlGit/queue-scheduler/internal/qs-api/api/v1/schedule"
	"github.com/SwirlGit/queue-scheduler/pkg/database/postgres"
	"github.com/SwirlGit/queue-scheduler/pkg/fasthttp"
)

const (
	appName        = "qs-api"
	configFilePath = "config.yaml"
)

func main() {
	cfg, err := config.InitConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	qsDB, err := postgres.NewDB(&cfg.QSDB, appName)
	if err != nil {
		panic(err)
	}

	scheduleStorage := pkgschedule.NewStorage(qsDB.Pool())
	scheduleService := schedule.NewService(scheduleStorage)
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
