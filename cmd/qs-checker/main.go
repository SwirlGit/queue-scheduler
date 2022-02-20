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
	"github.com/SwirlGit/queue-scheduler/pkg/log"
	"go.uber.org/zap"
)

const (
	appName           = "qs-checker"
	configFilePathEnv = "CONFIG_PATH"
)

func main() {
	cfg, err := config.InitConfig(configFilePathEnv)
	if err != nil {
		panic(err)
	}

	logger := log.NewZap(appName, zap.DebugLevel)
	defer func() { _ = logger.Sync() }()

	logger.Info("starting...")
	defer logger.Info("stopped")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	qsDB, err := postgres.NewDB(&cfg.QSDB, appName)
	if err != nil {
		logger.Panic("failed to init qs db", zap.Error(err))
	}

	scheduleStorage := pkgschedule.NewStorage(qsDB.Pool())
	checkerService := checker.NewService(logger, scheduleStorage, cfg.CheckDuration)
	checkerService.Start()
	defer checkerService.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("started")
	<-stop
	logger.Info("stopping...")
	signal.Stop(stop)
	close(stop)
}
