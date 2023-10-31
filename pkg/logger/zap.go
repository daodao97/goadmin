package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(zapLogger *zap.Logger) Logger {
	return &Zap{
		logger: zapLogger,
	}
}

type Zap struct {
	logger *zap.Logger
}

func (z Zap) Log(level Level, message string, keyValues ...interface{}) {
	var fields []zap.Field
	for index, v := range keyValues {
		if index%2 == 0 {
			fields = append(fields, zap.Any(v.(string), keyValues[index+1]))
		}
	}
	z.logger.Log(zapcore.Level(level), message, fields...)
}
