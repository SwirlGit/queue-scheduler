package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwirlGit/queue-scheduler/cmd/qs-worker/config"
	pkgschedule "github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/SwirlGit/queue-scheduler/internal/qs-worker/schedule"
	"github.com/SwirlGit/queue-scheduler/pkg/database/postgres"
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
	if err = scheduleService.Start(cfg.WorkersAmount); err != nil {
		panic(err)
	}
	defer scheduleService.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	signal.Stop(stop)
	close(stop)
}
