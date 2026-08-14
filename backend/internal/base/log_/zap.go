package log_

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

type ZapConfig struct {
	Mode string
}

func ToDefer() {
	err := Logger.Sync()
	if err != nil {
		return
	}
}

func InitZap(config ZapConfig) *zap.Logger {
	if config.Mode == "dev" {
		logger, _ := zap.NewDevelopment()
		Logger = logger
	} else {
		logger, _ := zap.NewProduction()
		Logger = logger
	}

	sugar := Logger.Sugar()
	Sugar = sugar
	return Logger
}

func Error(args ...interface{}) {
	Sugar.Error(args...)
}
func Info(args ...interface{}) {
	Sugar.Info(args...)
}

func Debug(args ...interface{}) {
	Sugar.Debug(args...)
}

func Warn(args ...interface{}) {
	Sugar.Warn(args...)
}
