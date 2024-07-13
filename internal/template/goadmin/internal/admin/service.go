package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/daodao97/goadmin/admin"
	"github.com/daodao97/goadmin/scaffold"
)

func ctrls() (routes []scaffold.GinRoute) {
	// 系统管理相关服务
	routes = admin.New()

	// 在此处注册业务的 Ctrl

	return routes
}

func NewRoute() (r []func(e *gin.Engine)) {
	for _, ctrl := range ctrls() {
		r = append(r, ctrl)
	}

	return r
}
