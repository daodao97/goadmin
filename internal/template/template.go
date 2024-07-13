package template

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/leaanthony/gosod"
)

//go:embed goadmin
var t embed.FS

type TmplVars struct {
	DriverMysql  bool
	DriverSqlite bool
	DefaultJWT   bool
}

func Install(projectPath string, data *TmplVars) error {
	defer func() {
		if data.DriverMysql {
			_ = os.Remove(filepath.Join(projectPath, "cmd", "goadmin.db"))
		}
	}()
	tfs, err := fs.Sub(t, "goadmin")
	if err != nil {
		return err
	}

	return gosod.New(tfs).Extract(projectPath, data)
}
