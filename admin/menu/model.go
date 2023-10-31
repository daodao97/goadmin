package menu

import (
	"fmt"
	"time"

	"github.com/daodao97/goadmin/pkg/db"
)

type menuReq struct {
	Id int64 `form:"id"`
}

type Menu struct {
	Page
	Children []*Menu `json:"children" form:"children"`
}

func (m *Menu) GetChildren(model db.Model) error {
	child := new([]*Menu)
	err := model.Select(db.WhereEq("pid", m.ID)).Binding(child)
	if err != nil {
		return err
	}
	if len(*child) == 0 {
		return fmt.Errorf("no children")
	}
	m.Children = *child
	return nil
}

const MenuDIRType = 1  // 目录
const MenuMenuType = 2 // 目录
const MenuPageType = 3 // 目录

type Page struct {
	ID         uint      `json:"id"          form:"id"`
	Pid        uint      `json:"pid"         form:"pid"`         // 父id
	ModuleID   int       `json:"module_id"   form:"module_id"`   // 模块id
	Name       string    `json:"name"        form:"name"`        // 路由名
	Type       uint8     `json:"type"        form:"type"`        // 类型  1 目录, 2 菜单, 3 页面
	Path       string    `json:"path"        form:"path"`        // 前端路由
	Icon       string    `json:"icon"        form:"icon"`        // 图标
	PageType   uint8     `json:"page_type"   form:"page_type"`   // 页面类型 0 自定义, 1 列表页, 2 表单页,  3复杂schema
	PageSchema string    `json:"page_schema" form:"page_schema"` // 页面定义
	View       string    `json:"view"        form:"view"`        // 自定义组件路径
	Sort       uint      `json:"sort"        form:"sort"`        // 倒序排序
	Status     int8      `json:"status"      form:"status"`      // 状态 0 禁用, 1 启用
	Ctime      time.Time `json:"ctime"`                          // 创建时间
	Mtime      time.Time `json:"mtime"`                          // 更新时间
}
