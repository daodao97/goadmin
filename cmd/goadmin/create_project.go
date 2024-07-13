package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/daodao97/goadmin/internal/template"
	"github.com/daodao97/goadmin/pkg/util"
)

type initArgs struct {
	ProjectName string `name:"n" description:"Name of project" default:""`
	Driver      string `name:"d" description:"db driver of project" default:"mysql"`
}

func createProject(options *initArgs) error {
	if options.ProjectName == "" {
		return fmt.Errorf("please use the -n flag to specify a project name")
	}

	vars := new(template.TmplVars)

	if options.Driver == "mysql" {
		vars.DriverMysql = true
	}
	if options.Driver == "sqlite" {
		vars.DriverSqlite = true
	}

	pwd, _ := os.Getwd()
	projectFullPath := filepath.Join(pwd, options.ProjectName)
	if util.DirectoryExists(projectFullPath) {
		return errors.New(fmt.Sprintf("Dir %s is exist", options.ProjectName))
	}

	err := template.Install(projectFullPath, vars)
	if err != nil {
		return err
	}

	fmt.Println("go mod tidy in " + projectFullPath)
	return util.StdCommand(context.Background(), projectFullPath, "go", "mod", "tidy")
}
