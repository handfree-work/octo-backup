package log_

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = zap.NewNop()
var Sugar = Logger.Sugar()
var logFile *lumberjack.Logger

type ZapConfig struct {
	Mode      string
	Directory string
	Level     string
}

func ToDefer() {
	_ = Logger.Sync()
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
	Logger = zap.NewNop()
	Sugar = Logger.Sugar()
}

// InitZap 同时输出到终端与滚动日志文件。
func InitZap(config ZapConfig) (*zap.Logger, error) {
	directory := config.Directory
	if directory == "" {
		directory = "./logs"
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, err
	}
	level := parseLevel(config.Level, config.Mode)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)
	logFile = &lumberjack.Logger{
		Filename:   filepath.Join(directory, "app.log"),
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		level,
	)
	Logger = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Sugar = Logger.Sugar()
	return Logger, nil
}

func HTTPMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		startedAt := time.Now()
		err := c.Next()
		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", time.Since(startedAt)),
		}
		if err != nil {
			Logger.Error("HTTP 请求失败", append(fields, zap.Error(err))...)
			return err
		}
		Logger.Info("HTTP 请求", fields...)
		return nil
	}
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

func parseLevel(value, mode string) zapcore.LevelEnabler {
	if value == "" && mode == "dev" {
		return zap.DebugLevel
	}
	switch strings.ToLower(value) {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
