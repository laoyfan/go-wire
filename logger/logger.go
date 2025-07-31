package logger

import (
	"context"
	"github.com/google/wire"
	"go-wire/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}

func ErrorField(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: "<nil>"}
	}
	return Field{Key: "error", Value: err.Error()}
}

var ProviderSet = wire.NewSet(NewLogger)

// NewLogger 根据配置返回对应的日志实现
func NewLogger(cfg *config.Config) (Logger, error) {
	switch cfg.Log.Driver {
	case "zap":
		return newZapLogger(cfg)
	default:
		// 默认使用 zap
		return newZapLogger(cfg)
	}
}

func newZapLogger(cfg *config.Config) (Logger, error) {
	logCfg := cfg.Log
	if err := ensureLogDirectoryExists(logCfg.Director); err != nil {
		return nil, err
	}
	writer := getLogWriter(logCfg.Director, logCfg.MaxSize, logCfg.MaxBackups, logCfg.MaxAge)

	// debug模式输出控制台
	if cfg.App.Mode == "debug" {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer)
	}
	// 创建编码器配置
	encoderConfig := getEncoderConfig()
	var core zapcore.Core
	if logCfg.Format == "json" {
		// 如果是JSON格式则使用JSONEncoder
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, levelPriority(getLevel(logCfg.Level)))
	} else {
		// 如果是Console格式则使用ConsoleEncoder
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, levelPriority(getLevel(logCfg.Level)))
	}

	log := zap.New(core)
	log = log.WithOptions(zap.AddCaller())

	return &zapLogger{logger: log}, nil
}
