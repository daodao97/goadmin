package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/daodao97/goadmin/internal/template"
	"github.com/daodao97/goadmin/pkg/util"
)

type startDemo struct{}

func startDemoProject(options *startDemo) error {
	tmpDir := os.TempDir()
	projectFullPath := filepath.Join(tmpDir, "goadmin_demo")
	if util.DirectoryExists(projectFullPath) {
		err := util.StdCommand(context.Background(), "", "rm", "-rf", projectFullPath)
		if err != nil {
			return err
		}
	}

	fmt.Println("Init Demo Project on", tmpDir)
	err := template.Install(projectFullPath, &template.TmplVars{
		DriverSqlite: true,
		DefaultJWT:   true,
	})
	if err != nil {
		return err
	}

	_, err = util.Command(projectFullPath, "go", "mod", "tidy")
	if err != nil {
		return err
	}

	fmt.Println("Start Project")
	os.Setenv("GIN_MODE", "release")
	return util.StdCommand(context.Background(), filepath.Join(projectFullPath, "cmd"), "go", "run", ".")
}
