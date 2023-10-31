package main

import (
	"fmt"
	"github.com/daodao97/goadmin/internal/gen"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type genArgs struct {
	Conn      string `name:"c" description:"Connection Name of db in config.Database" default:"default"`
	Database  string `name:"d" description:"Database Name" default:""`
	TableName string `name:"t" description:"Table Name"`
}

func genService(options *genArgs) error {
	if options.TableName == "" {
		return fmt.Errorf("please use the -t flag to set tableName")
	}

	pwd, _ := os.Getwd()
	projectFullPath := filepath.Join(pwd, options.TableName)
	if util.DirectoryExists(projectFullPath) {
		return errors.New(fmt.Sprintf("Dir %s is exist", options.TableName))
	}

	err := gen.CratePage(options.TableName, options.Conn, options.Database, options.TableName)
	if err != nil {
		return err
	}

	vars := &gen.TmplVars{
		PkgName:    options.TableName,
		PathPrefix: options.TableName,
		TableName:  options.TableName,
	}

	err = gen.Install(projectFullPath, vars)
	if err != nil {
		return errors.Errorf("install error: %v", err)
	}

	return nil
}
