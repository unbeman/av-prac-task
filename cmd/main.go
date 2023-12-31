package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-prac-task/internal/app"
	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/logging"
)

// @title Dynamic user segments server
// @version 1.0
// @description Avito homework.
// @BasePath /api/v1
func main() {
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

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		<-c
		sapp.Stop()
	}()

	sapp.Run()
}
