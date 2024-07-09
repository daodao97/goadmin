package scaffold

import (
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

// Ctrl 控制器接口
type Ctrl interface {
	Route(e *gin.Engine)   // 路由注册入口
	List(ctx *gin.Context) // 列表接口
	Get(ctx *gin.Context)  // 单一对象获取接口
	Set(ctx *gin.Context)  // 保存接口
	Del(ctx *gin.Context)  // 删除接口
}

type Server interface {
	Start() error
	Stop()
}

// BaseCtrl 当不需要实现脚手架增删改查接口时可用此结构体站位
type BaseCtrl struct{}

func (c *BaseCtrl) Route(e *gin.Engine)   {}
func (c *BaseCtrl) Get(ctx *gin.Context)  {}
func (c *BaseCtrl) Set(ctx *gin.Context)  {}
func (c *BaseCtrl) Del(ctx *gin.Context)  {}
func (c *BaseCtrl) List(ctx *gin.Context) {}

func NewFullCtrl(s *Scaffold, routeGroup string) FullCtrl {
	return FullCtrl{service: s, routeGroup: routeGroup}
}

// FullCtrl 接收设置好参数的 scaffold 补全基础的增删改查接口
type FullCtrl struct {
	routeGroup string
	service    *Scaffold
}

func (c *FullCtrl) Conf() *Conf {
	return c.service.Conf
}

func (c *FullCtrl) RegRoute(e *gin.Engine) *gin.RouterGroup {
	return RouteGroup(e, c.routeGroup, c, c.Conf(), func(ctx *gin.Context) {
		if c.routeGroup != ":table_name" {
			ctx.Next()
			return
		}
		tableName := ctx.Param("table_name")

		c.service.SetModel(db.New(tableName))

		ctx.Next()
	})
}

func (c *FullCtrl) Route(e *gin.Engine) {
	c.RegRoute(e)
}

func (c *FullCtrl) Get(ctx *gin.Context) {
	res, err := c.service.Get(ctx)
	Response(ctx, res, err)
}

func (c *FullCtrl) Set(ctx *gin.Context) {
	res, err := c.service.Set(ctx)
	Response(ctx, res, err)
}

func (c *FullCtrl) Del(ctx *gin.Context) {
	res, err := c.service.Del(ctx)
	Response(ctx, res, err)
}

func (c *FullCtrl) List(ctx *gin.Context) {
	res, err := c.service.List(ctx)
	Response(ctx, res, err)
}

func RouteGroup(_gin *gin.Engine, prefix string, ctrl Ctrl, conf *Conf, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	_prefix := prefix
	g := _gin.Group(filepath.Join(conf.HttpServer.BasePath, _prefix), handlers...)
	g.GET("/list", ctrl.List)       // 列表接口
	g.POST("/create", ctrl.Set)     // 创建接口
	g.GET("/get/:id", ctrl.Get)     // 单条获取
	g.POST("/update/:id", ctrl.Set) // 单条更新
	g.DELETE("/del", ctrl.Del)      // 单条/多条 删除 ?id=xx,xx

	return g
}
