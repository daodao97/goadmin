package user

import (
	"github.com/daodao97/goadmin/pkg/log"
	"github.com/daodao97/goadmin/scaffold/dao"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
	"strings"

	"github.com/daodao97/goadmin/scaffold"
)

func NewCtrl(s scaffold.Scaffold) scaffold.Ctrl {
	ser := newService(&s)
	return &ctrl{
		service:  ser,
		FullCtrl: scaffold.NewFullCtrl(ser.GetScaffold(), "user"),
	}
}

type ctrl struct {
	service *service
	scaffold.FullCtrl
}

func (c *ctrl) Route(e *gin.Engine) {
	g := c.FullCtrl.RegRoute(e)
	g.GET("/info", c.Info)
	g.GET("/routes", c.Routes)
	g.POST("/login", c.Login)
	g.GET("/logout", c.Logout)
	g.POST("/offline/:id", c.Offline)
	g.POST("/newpwd/:id", c.UpdatePwd)
	g.GET("/routesSE", c.RoutesSE)
	g.GET("/captcha", c.Captcha)
	g.GET("/select_options", c.selectOption)
}

// Info 当前登录用户的信息获取接口
func (c *ctrl) Info(ctx *gin.Context) {
	res, err := c.service.User(ctx)
	scaffold.Response(ctx, res, err)
}

// Routes 当前登录用户有权限的路由
func (c *ctrl) Routes(ctx *gin.Context) {
	res, err := c.service.Routes(ctx)
	scaffold.Response(ctx, res, err)
}

// Login 登录接口
func (c *ctrl) Login(ctx *gin.Context) {
	res, err := c.service.Login(ctx)
	scaffold.Response(ctx, res, err)
}

// Logout 登出接口
func (c *ctrl) Logout(ctx *gin.Context) {
	err := c.service.Logout(ctx)
	scaffold.Response(ctx, nil, err)
}

// Offline 下线接口
func (c *ctrl) Offline(ctx *gin.Context) {
	err := c.service.Offline(ctx)
	scaffold.Response(ctx, nil, err)
}

// UpdatePwd 改密接口
func (c *ctrl) UpdatePwd(ctx *gin.Context) {
	err := c.service.UpdatePwd(ctx)
	scaffold.Response(ctx, nil, err)
}

// RoutesSE 当前登录用户有权限的路由精简映射
func (c *ctrl) RoutesSE(ctx *gin.Context) {
	res, err := c.service.RoutesSE(ctx)
	scaffold.Response(ctx, res, err)
}

func (c *ctrl) Captcha(ctx *gin.Context) {
	captchaText := GetRandStr(4)
	if err := c.service.UserState.SetCaptcha(ctx, strings.ToLower(captchaText)); err != nil {
		log.Error("Captcha err(%+v)", err) // login校验不会通过了
		return
	}
	ctx.Render(http.StatusOK, render.Data{
		ContentType: "image/png",
		Data:        CaptchaImgText(180, 60, captchaText),
	})
}

func (c *ctrl) selectOption(ctx *gin.Context) {
	req := new(struct {
		Kw string `form:"kw"`
	})
	if err := ctx.Bind(req); err != nil {
		scaffold.Response(ctx, nil, err)
		return
	}
	res, err := c.service.SelectOption(ctx, req.Kw, "id", "nickname", dao.Field("id", "nickname"))
	scaffold.Response(ctx, res, err)
}
