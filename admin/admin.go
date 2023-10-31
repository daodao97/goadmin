package admin

import (
	"github.com/daodao97/goadmin/admin/cconf"
	"github.com/daodao97/goadmin/admin/menu"
	"github.com/daodao97/goadmin/admin/role"
	"github.com/daodao97/goadmin/admin/user"
	"github.com/daodao97/goadmin/scaffold"
)

func New(s scaffold.Scaffold) []scaffold.Ctrl {
	ctrl := []scaffold.Ctrl{
		role.NewCtrl(s),
		user.NewCtrl(s),
		menu.NewCtrl(s),
		cconf.NewCtrl(s),
	}
	return ctrl
}
