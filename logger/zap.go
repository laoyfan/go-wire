package logger

import (
	"context"
	"fmt"
	"go-wire/config"
	"go-wire/util"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	log *zap.Logger
}

func NewZapLogger(cfg *config.Config) (Logger, error) {
	if err := ensureLogDirectoryExists(cfg.Log.Director); err != nil {
		return nil, err
	}
	writer := getLogWriter(cfg.Log.Director, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)

	// debug模式输出控制台
	if cfg.App.Mode == "debug" {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer)
	}
	// 创建编码器配置
	encoderConfig := getEncoderConfig()
	var core zapcore.Core
	if cfg.Log.Format == "json" {
		// 如果是JSON格式则使用JSONEncoder
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, levelPriority(getLevel(cfg.Log.Level)))
	} else {
		// 如果是Console格式则使用ConsoleEncoder
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, levelPriority(getLevel(cfg.Log.Level)))
	}

	log := zap.New(core)
	log = log.WithOptions(zap.AddCaller())

	return &zapLogger{log: log}, nil
}

// FieldErr 用于将 error 包装成 zap.Field，方便统一日志格式。
func (l *zapLogger) FieldErr(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

// 转换自定义字段到 zap.Field
func (l *zapLogger) toZapFields(fields []Field) []zap.Field {
	zFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zFields = append(zFields, zap.Any(f.Key, f.Value))
	}
	return zFields
}

// logWithTraceID 带有 TraceID 的日志记录
func (l *zapLogger) withTrace(ctx context.Context, level zapcore.Level, msg string, fields ...Field) {
	traceID, _ := ctx.Value("TraceID").(string)
	if l.log.Core().Enabled(level) {
		l.log.With(zap.Any("trace_id", traceID)).WithOptions(zap.AddCallerSkip(2)).Log(level, msg, l.toZapFields(fields)...)
	}
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.withTrace(ctx, zapcore.DebugLevel, msg, fields...)
}
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.withTrace(ctx, zapcore.InfoLevel, msg, fields...)
}
func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.withTrace(ctx, zapcore.WarnLevel, msg, fields...)
}
func (l *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.withTrace(ctx, zapcore.ErrorLevel, msg, fields...)
}
func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.withTrace(ctx, zapcore.FatalLevel, msg, fields...)
}

func (l *zapLogger) Sync() error {
	return l.log.Sync()
}

// 创建日志目录
func ensureLogDirectoryExists(director string) error {
	if ok, _ := util.PathExists(director); !ok {
		fmt.Println("创建日志文件夹", director)
		if err := os.Mkdir(director, os.ModePerm); err != nil {
			return fmt.Errorf("创建日志文件夹失败: %w", err)
		}
	}
	return nil
}

// 获取日志文件写入器
func getLogWriter(director string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	logFileName := path.Join(director, "service.log")
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})
}

// 获取编码配置
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// 获取日志级别
func getLevel(level string) zapcore.Level {
	levelMap := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	if zapLevel, exists := levelMap[level]; exists {
		return zapLevel
	}
	return zapcore.DebugLevel
}

// 获取日志级别优先级
func levelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool {
		return l >= level
	}
}
