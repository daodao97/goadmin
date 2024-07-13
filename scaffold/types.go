package scaffold

import (
	"encoding/json"

	"github.com/daodao97/goadmin/scaffold/dao"
)

// 脚手架的通用类型定义

type Page struct {
	Pn    int64 `json:"pn"`
	Ps    int64 `json:"ps"`
	Total int64 `json:"total"`
}

type ListResp struct {
	Page Page      `json:"page"`
	List []dao.Row `json:"list"`
}

type ListRespInterface struct {
	Page Page        `json:"page"`
	List interface{} `json:"list"`
}

type TreeList struct {
	PidKey string `json:"pid_key"`
	Lazy   bool   `json:"lazy"`
}

func (t *TreeList) UnmarshalJSON(data []byte) error {
	t.PidKey = "id"
	type Alias TreeList
	tmp := (*Alias)(t)

	return json.Unmarshal(data, tmp)
}

const (
	PageTypeCustom = 0
	PageTypeList   = 1
	PageTypeForm   = 2
	PageTypeSchema = 3
)

type Route struct {
	ID         uint        `json:"id"`
	Pid        uint        `json:"pid"`         // 父id
	ModuleID   int         `json:"module_id"`   // 模块id
	Name       string      `json:"name"`        // 路由名
	Type       uint8       `json:"type"`        // 类型  1 目录, 2 菜单, 3 页面
	Path       string      `json:"path"`        // 前端路由
	Icon       string      `json:"icon"`        // 图标
	PageType   uint8       `json:"page_type"`   // 页面类型 0 自定义, 1 列表页, 2 表单页,  3复杂schema
	PageSchema interface{} `json:"page_schema"` // 页面定义
	View       string      `json:"view"`        // 自定义组件路径
	Sort       uint        `json:"sort"`        // 倒序排序
	Status     int8        `json:"status"`      // 状态 0 禁用, 1 启用
	Children   []Route     `json:"children"`
}

type Module struct {
	Id     int     `json:"id"`
	Label  string  `json:"label"`
	Routes []Route `json:"routes"`
}

type Website struct {
	Title                string             `json:"title,omitempty"`
	FixedHeader          bool               `json:"fixedHeader,omitempty"`
	SidebarLogo          bool               `json:"sidebarLogo,omitempty"`
	Logo                 string             `json:"logo,omitempty"`
	CloseNavNotice       bool               `json:"closeNavNotice,omitempty"`
	NavBarNotice         string             `json:"navBarNotice,omitempty"`
	HasNewMessage        bool               `json:"hasNewMessage,omitempty"`
	ShowPageJsonSchema   bool               `json:"showPageJsonSchema,omitempty"`
	LoginTips            string             `json:"loginTips,omitempty"`
	ElementPlus          *ElementPlusConfig `json:"ElementPlus,omitempty"`
	WhiteRoutes          []string           `json:"whiteRoutes,omitempty"`
	TokenExpire          int64              `json:"tokenExpire,omitempty"`
	DefaultAvatar        string             `json:"defaultAvatar,omitempty"`
	EnvColor             map[string]string  `json:"envColor,omitempty"`
	ServiceOfflineNotice string             `json:"serviceOffLineNotice,omitempty"`
	Modules              []Module           `json:"modules,omitempty"`
	MacroVar             string             `json:"macroVar,omitempty"`
}

type ElementPlusConfig struct {
	Size   string `json:"size,omitempty"`
	ZIndex int    `json:"zIndex,omitempty"`
	Locale string `json:"locale,omitempty"`
}

type UserAttr struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"     comment:"名称"`
	Nickname string          `json:"nickname" comment:"名称"`
	Avatar   string          `json:"avatar"   comment:"头像"`
	Email    string          `json:"email"    comment:"邮箱"`
	Mobile   string          `json:"mobile"   comment:"手机号"`
	RoleIds  [][]int         `json:"role_ids" comment:"用户角色"`
	Resource [][]interface{} `json:"resource"`
	Env      string          `json:"env"`
	Website  *Website        `json:"website"`
}

func (u *UserAttr) IsSupper() bool {
	const id = 1
	for _, v := range u.RoleIds {
		if v[len(v)-1] == id {
			return true
		}
	}

	return false
}
