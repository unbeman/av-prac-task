package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/unbeman/av-prac-task/internal/app"
	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/logging"
)

// @title Dynamic user segments server
// @version 1.0
// @description Avito homework.
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		<-c
		cancel()
	}()

	cfg, err := config.GetAppConfig()
	if err != nil {
		log.Error("can't get config: ", err)
		return
	}

	logging.InitLogger(cfg.Logger)
	sapp, err := app.GetSegApp(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	sapp.Run(ctx)
}
