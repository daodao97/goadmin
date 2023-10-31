package service

import (
	"github.com/google/wire"

	"admin/internal/app"
	"github.com/daodao97/goadmin/admin"
	"github.com/daodao97/goadmin/scaffold"
)

var Provider = wire.NewSet(NewCtrl, wire.Struct(new(CtrlOptions), "*"))

type CtrlOptions struct {
	Scaffold scaffold.Scaffold
	Conf     *app.Conf
}

// NewCtrl 向外导出所有的 http路由 服务
func NewCtrl(opt *CtrlOptions) (ctrl []scaffold.Ctrl, cf func(), err error) {
	// 系统管理相关服务
	ctrl = admin.New(opt.Scaffold)
	// 在此处注册业务的 Ctrl

	cf = func() {}
	return ctrl, cf, nil
}
