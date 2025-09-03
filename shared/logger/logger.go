package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/Prabhat7saini/Basic-Setup/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logInstance *zap.Logger
	once        sync.Once
)

// Singleton zap logger initialization
func Init(cfg config.LogConfig) *zap.Logger {
	once.Do(func() {
		logInstance = buildLogger(cfg)
	})
	return logInstance
}

// builds logger based on config
func buildLogger(cfg config.LogConfig) *zap.Logger {
	// Parse level
	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	// encoder
	var encoder zapcore.Encoder
	if strings.ToLower(cfg.Format) == "json" {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lvl)

	var options []zap.Option
	if cfg.EnableCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger := zap.New(core, options...)
	return logger
}

// returns the singleton logger
func Get() *zap.Logger {
	if logInstance == nil {
		return buildLogger(config.LogConfig{
			Level:            "info",
			Format:           "console",
			EnableCaller:     true,
			EnableStacktrace: false,
		})
	}
	return logInstance
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func WithContext(ctx context.Context) *zap.Logger {
	if requestID := ctx.Value("request_id"); requestID != nil {
		return Get().With(zap.String("request_id", requestID.(string)))
	}
	return Get()
}
