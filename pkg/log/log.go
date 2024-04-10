package log

import (
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var opts = slog.HandlerOptions{
	Level: slog.LevelDebug,
}

var logger = StdoutText()

func SetLogger(l *slog.Logger) {
	logger = l
	slog.SetDefault(l)
}

func StdoutText() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &opts))
}

func StdoutTextPretty() *slog.Logger {
	return slog.New(NewPrettyHandler(os.Stdout, PrettyHandlerOptions{
		SlogOpts: opts,
	}))
}

func StdoutJson() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &opts))
}

func FileJson(fileName string) *slog.Logger {
	r := &lumberjack.Logger{
		Filename:   fileName,
		LocalTime:  true,
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 5,
		Compress:   true,
	}
	return slog.New(slog.NewJSONHandler(r, &opts))
}

func Debug(msg string, kv ...any) {
	logger.Debug(msg, kv...)
}

func Info(msg string, kv ...any) {
	logger.Info(msg, kv...)
}

func Error(msg string, kv ...any) {
	logger.Error(msg, kv...)
}

func Warn(msg string, kv ...any) {
	logger.Warn(msg, kv...)
}
