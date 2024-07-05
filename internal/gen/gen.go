package gen

import (
	"embed"
	"github.com/leaanthony/gosod"
	"io/fs"
)

//go:embed service
var t embed.FS

type TmplVars struct {
	PkgName    string
	PathPrefix string
	TableName  string
}

func Install(projectPath string, data *TmplVars) error {
	tfs, err := fs.Sub(t, "admin")
	if err != nil {
		return err
	}

	return gosod.New(tfs).Extract(projectPath, data)
}
