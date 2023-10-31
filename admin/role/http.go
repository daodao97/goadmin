package role

import (
	"github.com/gin-gonic/gin"

	"github.com/daodao97/goadmin/scaffold"
)

func NewCtrl(s scaffold.Scaffold) scaffold.Ctrl {
	ser := newService(&s)
	return &ctrl{
		service:  ser,
		FullCtrl: scaffold.NewFullCtrl(ser.GetScaffold(), "role"),
	}
}

type ctrl struct {
	service *service
	scaffold.FullCtrl
}

func (c *ctrl) Route(e *gin.Engine) {
	g := c.FullCtrl.RegRoute(e)
	g.GET("/tree", c.Tree)
	g.GET("/resource", c.Resource)
}

func (c *ctrl) Tree(ctx *gin.Context) {
	res, err := c.service.Tree(ctx)
	scaffold.Response(ctx, res, err)
}

func (c *ctrl) Resource(ctx *gin.Context) {
	res, err := c.service.Resource(ctx)
	scaffold.Response(ctx, res, err)
}
