package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwirlGit/queue-scheduler/cmd/qs-api/config"
	pkgschedule "github.com/SwirlGit/queue-scheduler/internal/pkg/schedule"
	"github.com/SwirlGit/queue-scheduler/internal/qs-api/api/v1/schedule"
	"github.com/SwirlGit/queue-scheduler/pkg/database/postgres"
	"github.com/SwirlGit/queue-scheduler/pkg/fasthttp"
	"github.com/SwirlGit/queue-scheduler/pkg/log"
	"go.uber.org/zap"
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

	logger := log.NewZap(appName, zap.DebugLevel)
	defer func() { _ = logger.Sync() }()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	qsDB, err := postgres.NewDB(&cfg.QSDB, appName)
	if err != nil {
		logger.Panic("failed to init qs db", zap.Error(err))
	}

	scheduleStorage := pkgschedule.NewStorage(qsDB.Pool())
	scheduleService := schedule.NewService(scheduleStorage)
	scheduleHandler := schedule.NewHandler(scheduleService)

	server := fasthttp.NewServer([]fasthttp.RouteProvider{scheduleHandler})
	go func() {
		if err := server.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			logger.Panic("failed to start listen", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	signal.Stop(stop)
	close(stop)

	if err := server.Shutdown(); err != nil {
		logger.Panic("failed to shutdown server", zap.Error(err))
	}
}
