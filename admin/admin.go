package admin

import (
	"github.com/daodao97/goadmin/admin/cconf"
	"github.com/daodao97/goadmin/admin/menu"
	"github.com/daodao97/goadmin/admin/role"
	"github.com/daodao97/goadmin/admin/user"
	"github.com/daodao97/goadmin/scaffold"
)

func New() (routes []scaffold.GinRoute) {
	s := scaffold.GetScaffold()
	ctrl := []scaffold.Ctrl{
		role.NewCtrl(s),
		user.NewCtrl(s),
		menu.NewCtrl(s),
		cconf.NewCtrl(s),
	}

	for _, r := range ctrl {
		routes = append(routes, r.Route)
	}

	return
}
