package logger

import (
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewZapLogger)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
	Sync() error
}

type Field struct {
	Key   string
	Value any
}

func KeyValue(key string, value any) Field {
	return Field{Key: key, Value: value}
}

func StringAny(key string, value any) Field {
	return KeyValue(key, value)
}

func Error(err error) Field {
	return KeyValue("error", err)
}
