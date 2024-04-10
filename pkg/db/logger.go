package db

import (
	"time"

	"github.com/daodao97/goadmin/pkg/log"
)

func Info(msg string, kv ...interface{}) {
	log.Info(msg, kv...)
}

func Error(msg string, kv ...interface{}) {
	log.Error(msg, kv...)
}

func dbLog(prefix string, start time.Time, err *error, kv *[]interface{}) {
	tc := time.Since(start)
	_log := []interface{}{
		"scope", "db",
		"prefix", prefix,
		"ums", tc.Milliseconds(),
	}
	_log = append(_log, *kv...)
	if *err != nil {
		_log = append(_log, "error", *err)
		log.Error("query", _log...)
		return
	}
	log.Debug("query", _log...)
}
