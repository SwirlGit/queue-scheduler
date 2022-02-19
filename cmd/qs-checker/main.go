package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwirlGit/queue-scheduler/cmd/qs-checker/config"
	pkgschedule "github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/SwirlGit/queue-scheduler/internal/qs-checker/checker"
	"github.com/SwirlGit/queue-scheduler/pkg/database/postgres"
)

const (
	appName        = "qs-checker"
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
	checkerService := checker.NewService(scheduleStorage, cfg.CheckDuration)
	checkerService.Start()
	defer checkerService.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	signal.Stop(stop)
	close(stop)
}
