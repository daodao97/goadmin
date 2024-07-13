package menu

import (
	"encoding/json"

	"github.com/daodao97/goadmin/pkg/ecode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/daodao97/goadmin/scaffold"
)

func NewCtrl(s scaffold.Scaffold) scaffold.Ctrl {
	ser := newService(&s)
	return &ctrl{
		service:  ser,
		FullCtrl: scaffold.NewFullCtrl(ser.GetScaffold(), "menu"),
	}
}

type ctrl struct {
	service *service
	scaffold.FullCtrl
}

func (c *ctrl) Route(e *gin.Engine) {
	g := c.FullCtrl.RegRoute(e)
	g.GET("/tree", c.Tree)
	g.GET("/children", c.Children)
	g.POST("/import", c.Import)
	g.GET("/export", c.Export)
	g.GET("/mm", func(ctx *gin.Context) {
		err := c.service.MenuMigration(ctx)
		scaffold.Response(ctx, nil, err)
	})
}

func (c *ctrl) Children(ctx *gin.Context) {
	res, err := c.service.Children(ctx)
	scaffold.Response(ctx, res, err)
}

func (c *ctrl) Tree(ctx *gin.Context) {
	res, err := c.service.Tree(ctx)
	scaffold.Response(ctx, res, err)
}

func (c *ctrl) Export(ctx *gin.Context) {
	req := new(menuReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	res, err := c.service.ExportMenu(ctx, req)
	scaffold.Response(ctx, map[string]interface{}{"data": res}, err)
}

func (c *ctrl) Import(ctx *gin.Context) {
	id := new(menuReq)
	if err := ctx.Bind(id); err != nil {
		return
	}
	menus := new(struct {
		Data string `form:"data" validate:"required"`
	})
	if err := ctx.BindWith(menus, binding.JSON); err != nil {
		return
	}
	menu := new(Menu)
	err := json.Unmarshal([]byte(menus.Data), menu)
	if err != nil {
		scaffold.RenderErrMsg(ctx, ecode.RequestErr.Code(), "导入数据解析失败")
		return
	}
	res, err := c.service.ImportMenu(ctx, id.Id, menu)
	scaffold.Response(ctx, map[string]bool{"status": res}, err)
}
