package cconf

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/ecode"
	"github.com/daodao97/goadmin/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
)

func NewCtrl(s scaffold.Scaffold) scaffold.Ctrl {
	ser := newService(&s)
	return &ctrl{
		service:  ser,
		FullCtrl: scaffold.NewFullCtrl(ser.GetScaffold(), "cconf"),
	}
}

type ctrl struct {
	service *service
	scaffold.FullCtrl
}

func (c *ctrl) Route(e *gin.Engine) {
	g := c.FullCtrl.RegRoute(e)
	// 表单填写页面接口
	g.GET("/form_schema/:id", c.FormSchema) // 表单 schema
	g.GET("/form_value/:id", c.FormValue)   // 表单现有值
	g.POST("/form_save/:id", c.FormSave)    // 保存保存
}

func (c *ctrl) FormSchema(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		scaffold.Response(ctx, nil, fmt.Errorf("缺少记录id"))
		return
	}
	res, err := c.service.getCommConfById(ctx, id)
	if err != nil {
		return
	}
	schema, ok := res["rules"]
	if !ok || schema == nil {
		scaffold.Response(ctx, nil, fmt.Errorf("表单配置不可用"))
		return
	}
	_schema, err := util.String(cast.ToString(schema)).DecodeInterfaceE()
	scaffold.Response(ctx, _schema, err)
}

func (c *ctrl) FormValue(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		scaffold.Response(ctx, nil, fmt.Errorf("缺少记录id"))
		return
	}
	res, err := c.service.getCommConfById(ctx, id)
	if err != nil {
		return
	}
	scaffold.Response(ctx, res["value"], err)
}

func (c *ctrl) FormSave(ctx *gin.Context) {
	id, exist := ctx.Params.Get("id")
	if !exist {
		scaffold.Response(ctx, nil, ecode.ServerErr.Message("获取记录必须有主键 如 /user/1"))
		return
	}
	row, err := c.service.getCommConfById(ctx, id)
	if err != nil {
		scaffold.Response(ctx, nil, ecode.ServerErr.Message(err.Error()))
		return
	}
	m := c.service.GetModel(ctx)
	if err != nil {
		scaffold.Response(ctx, nil, ecode.ServerErr.Message(err.Error()))
		return
	}
	body, err := c.service.RequestBody(ctx)
	if err != nil {
		scaffold.Response(ctx, nil, ecode.ServerErr.Message(err.Error()))
		return
	}
	result, err := m.Update(map[string]interface{}{"value": body}, db.WhereEq("id", id))
	if err != nil {
		scaffold.Response(ctx, nil, ecode.ServerErr.Message(err.Error()))
		return
	}
	_ = c.service.CommonConf.Delete(ctx, cast.ToString(row["name"]))
	scaffold.Response(ctx, result, nil)
}
