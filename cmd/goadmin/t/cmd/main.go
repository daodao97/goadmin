package main

import (
	"flag"
	"time"

	"github.com/daodao97/xgo/xlog"
	_ "github.com/go-sql-driver/mysql"
	_ "go.uber.org/automaxprocs"
)

func main() {
	flag.Parse()
	app, closeFunc, err := initApp()
	if err != nil {
		panic(err)
	}

	exit, err := app.Start()
	if err != nil {
		panic(err)
	}

	for range exit {
		closeFunc()
		xlog.Info("admin exit")
		time.Sleep(time.Second)
		return
	}
}
