package app

import (
	"admin/internal/conf"
	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/sso"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/wire"
	"github.com/jinzhu/copier"

	"github.com/daodao97/goadmin/pkg/util/uploader"
	"github.com/daodao97/goadmin/scaffold"
)

var Provider = wire.NewSet(
	NewHttpServer, NewApp, wire.Struct(new(HttpOptions), "*"),
	NewCache, NewUploader, NewScaffoldConf, NewSso,
)

func NewApp(hs *HttpServer) *App {
	return &App{
		http: hs,
	}
}

type App struct {
	http *HttpServer
}

func (a *App) Start() (<-chan os.Signal, error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	return c, nil
}

func NewCache(c *conf.Conf) cache.Cache {
	return cache.NewRedisCache(c.Redis)
}

func NewUploader(c *conf.Conf) uploader.Uploader {
	// 根据情况选择不同的 Uploader
	return uploader.NewLocalUploader(
		"./uploads",
		"ok/{year}{month}{day}-{hour}{minute}{second}-{random}{.suffix}",
		"http://127.0.0.1:8001",
	)
}

func NewScaffoldConf(c *conf.Conf) (*scaffold.Conf, error) {
	cf := &scaffold.Conf{}
	err := copier.Copy(cf, c)
	if err != nil {
		return nil, err
	}
	err = cf.Validate()
	if err != nil {
		return nil, err
	}

	return cf, nil
}

func NewSso() *sso.Sso {
	// 若没有统一登录可返回空map
	return &map[sso.Name]sso.SSO{}
}
