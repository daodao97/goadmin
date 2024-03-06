package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).Sprint
	yellow = color.New(color.FgYellow).Sprint
	green  = color.New(color.FgGreen).Sprint
	gray   = color.New(color.FgBlue).Sprint
	prefix = color.New(color.FgGreen).Sprint
)

func newStdOutLogger() Logger {
	return &stdOutLogger{
		logger: log.New(os.Stdout, prefix("LOG: "), log.LstdFlags),
	}
}

type stdOutLogger struct {
	logger *log.Logger
}

func jsonEncode(v interface{}) string {
	bt, _ := json.Marshal(v)
	return string(bt)
}

func (s stdOutLogger) Log(level Level, message string, keyValues ...interface{}) {
	if level < limitLevel {
		return
	}
	args := []interface{}{message}
	if withCaller {
		args = append(args, "caller", caller(5))
	}
	args = append(args, keyValues...)

	var _args []interface{}
	_args = append(_args, colorLevel(level))
	for _, v := range args {
		switch t := v.(type) {
		case []interface{}:
			_args = append(_args, jsonEncode(t))
		default:
			_args = append(_args, v)
		}
	}

	s.logger.Println(_args...)
}

func caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}

func colorLevel(level Level) string {
	fn := gray
	switch level {
	case LevelInfo:
		fn = green
	case LevelWarn:
		fn = yellow
	case LevelError:
		fn = red
	}
	return fn(level.String())
}
