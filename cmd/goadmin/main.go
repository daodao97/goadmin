package main

import (
	"github.com/leaanthony/clir"
	"github.com/pterm/pterm"
	"os"
)

func main() {
	app := clir.NewCli("GoAdmin", "The GoAdmin CLI", "v0.0.1")
	app.NewSubCommandFunction("create", "Create a new GoAdmin project", createProject)
	app.NewSubCommandFunction("start_demo", "Start GoAdmin demo", startDemoProject)
	app.NewSubCommandFunction("gen", "Generate admin of project", genService)
	err := app.Run()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
