package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/goProjects/loan_app/lib/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger      *zap.Logger
	jsonEncoder zapcore.Encoder
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Init(mode int, env utils.Environment) {
	var logLevel zapcore.Level
	switch mode {
	case DEBUG:
		logLevel = zapcore.DebugLevel
	case INFO:
		logLevel = zapcore.InfoLevel
	case WARNING:
		logLevel = zapcore.WarnLevel
	case ERROR:
		logLevel = zapcore.ErrorLevel
	case FATAL:
		logLevel = zapcore.FatalLevel
	}

	cfg := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(logLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, _ = cfg.Build()
	jsonEncoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)

	if env == utils.DevEnv {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = logger.WithOptions(
			zap.WrapCore(
				func(zapcore.Core) zapcore.Core {
					return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.AddSync(os.Stderr), zapcore.DebugLevel)
				}))
	} else {
		logger = logger.WithOptions(
			zap.WrapCore(
				func(zapcore.Core) zapcore.Core {
					return zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stderr), logLevel)
				}))
	}
}

func Get() *zap.Logger {
	return logger
}

func Field(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

func I(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Info(message, addRequestID(ctx, fields...)...)
}

func D(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Debug(message, addRequestID(ctx, fields...)...)
}

func W(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Warn(message, addRequestID(ctx, fields...)...)
}

func E(ctx context.Context, message string, fields ...zapcore.Field) {
	fieldsWithRequestID := addRequestID(ctx, fields...)
	logger.Error(message, fieldsWithRequestID...)
	// notify error on some external tool like airbrake
}

func Sync() {
	logger.Info("SYNCING LOGGER....")
	if err := logger.Sync(); err != nil {
		fmt.Println("FAILED TO SYNC LOGGER...")
	}
}

func addRequestID(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	if requestID := utils.GetRequestID(ctx); requestID != "" {
		fields = append(fields, zap.String(string(utils.RequestIDLogKey), requestID))
	}
	return fields
}
