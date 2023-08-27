package logging

import (
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-prac-task/internal/config"
)

const (
	LogDebug = "debug"
	LogInfo  = "info"
)

func InitLogger(cfg config.LoggerConfig) {
	switch cfg.Level {
	case LogInfo:
		log.SetLevel(log.InfoLevel)
	case LogDebug:
		log.SetLevel(log.DebugLevel)
	}
}
