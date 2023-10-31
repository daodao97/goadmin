package logger

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_stdout(t *testing.T) {
	_logger.Log(LevelError, "stdout", "test", 222)
	_logger.Log(LevelInfo, "stdout", " test", 222)
	_logger.Log(LevelDebug, "stdout", "test", 222)
}

func Test_zap(t *testing.T) {
	log, _ := zap.NewDevelopment()
	zapLog := NewZap(log)
	SetLoggerLevel(LevelDebug)
	SetLogger(zapLog)
	_logger.Log(LevelDebug, "zap log", "test", 222, "tt", time.Now())
}
