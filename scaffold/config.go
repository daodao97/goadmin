package scaffold

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/xgo/xlog"
	"github.com/pkg/errors"
	"time"

	"github.com/daodao97/goadmin/pkg/db"
	"path/filepath"
)

type DBConf struct {
	Name string `validate:"required"`
	db.Config
}

type Jwt struct {
	Secret      string `validate:"required"`
	TokenExpire int64  `validate:"required,gt=3600"`
	OpenApi     []string
	PublicApi   []string
}

type HttpServer struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	BasePath       string
	WebPath        string
}

// Conf 脚手架依赖的核心配置
type Conf struct {
	HttpServer *HttpServer
	Database   []*DBConf `validate:"required"`
	Jwt        *Jwt      `validate:"required"`
}

func (c Conf) Validate() error {
	errTpl := "scaffold config check : %s"
	if c.HttpServer == nil {
		return fmt.Errorf(errTpl, "HttpServer is required")
	}
	if c.HttpServer.BasePath == "" {
		c.HttpServer.BasePath = "/_api"
	}

	if c.HttpServer.WebPath != "" {
		if c.HttpServer.WebPath == "/" {
			return fmt.Errorf(errTpl, "HttpServer.webPath can not be start with `/`")
		}
		if !util.String(c.HttpServer.WebPath).StartWith("/") {
			return fmt.Errorf(errTpl, "HttpServer.webPath nust start with `/` ")
		}
	}

	if c.Jwt == nil {
		return fmt.Errorf(errTpl, "Jwt is required")
	}
	if c.Jwt.Secret == "{token_secret}" {
		xlog.Warn("Place change Jwt.Secret = {token_secret} in your application.toml")
	}

	c.Jwt.OpenApi = append(c.Jwt.OpenApi, []string{
		filepath.Join(c.HttpServer.BasePath, "/user/login"),
		filepath.Join(c.HttpServer.BasePath, "/user/captcha*"),
	}...)
	c.Jwt.PublicApi = append(c.Jwt.OpenApi, []string{
		filepath.Join(c.HttpServer.BasePath, "/user/info"),
		filepath.Join(c.HttpServer.BasePath, "/user/routes"),
	}...)
	if len(c.Database) == 0 {
		return fmt.Errorf(errTpl, "Database is required")
	}
	hasDefaultDb := false
	for _, v := range c.Database {
		if v.Name == "default" {
			hasDefaultDb = true
		}
	}
	if !hasDefaultDb {
		return fmt.Errorf(errTpl, "Must have one database config with `name=default`")
	}
	err := util.NewValidate().Struct(c)
	if err != nil {
		return fmt.Errorf(errTpl, err.Error())
	}
	return nil
}

func (c Conf) GetDefaultDB() (*DBConf, error) {
	for _, v := range c.Database {
		if v.Name == "default" {
			return v, nil
		}
	}
	return nil, errors.New("can not find default database")
}

func (c Conf) GetDBMap() map[string]*db.Config {
	dbmap := make(map[string]*db.Config)
	for _, v := range c.Database {
		dbmap[v.Name] = &v.Config
	}
	return dbmap
}
