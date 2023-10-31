package menu

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"

	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
	"github.com/daodao97/goadmin/scaffold/dao"
	"github.com/gin-gonic/gin"
)

func newService(s *scaffold.Scaffold) *service {
	s.SetModel(db.New("page", db.ColumnHook(db.CommaInt("owner_ids"))))
	s.TreeList = &scaffold.TreeList{Lazy: true}
	s.BeforeCreate = func(ctx *gin.Context, val dao.Row) (dao.Row, error) {
		if val["path"] == "#" {
			return val, nil
		}
		count, err := s.GetModel(ctx).Count(db.WhereEq("path", val["path"]))
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("该路由已存在 %s", val["path"])
		}
		return val, nil
	}
	s.BeforeUpdate = func(ctx *gin.Context, val dao.Row, id int64) (dao.Row, error) {
		if val["path"] == "#" {
			return val, nil
		}
		u, err := s.User(ctx)
		if err != nil {
			return nil, err
		}

		op := func(id int, ids []int) bool {
			if u.IsSupper() {
				return true
			}
			for _, v := range ids {
				if v == id {
					return true
				}
			}
			return false
		}(u.Id, cast.ToIntSlice(val["owner_ids"]))
		if !op {
			return nil, errors.New("您没有权限更改此页面的配置")
		}
		row := s.GetModel(ctx).SelectOne(db.WhereEq("path", val["path"]))
		if row.Err != nil && !errors.Is(row.Err, db.ErrNotFound) {
			return nil, row.Err
		}
		if row.Data == nil || cast.ToInt64(row.Data["id"]) == id {
			return val, nil
		}
		return nil, fmt.Errorf("该路由已存在 %s", val["path"])
	}
	return &service{
		Scaffold: s,
	}
}

type service struct {
	*scaffold.Scaffold
}

func (s *service) GetScaffold() *scaffold.Scaffold {
	return s.Scaffold
}

func (s *service) Tree(ctx *gin.Context) ([]util.TreeData, error) {
	m := s.GetModel(ctx)
	rows := m.Select(db.Field("id", "pid", "name"), db.WhereLt("type", 3))
	if rows.Err != nil {
		return nil, rows.Err
	}
	var tree []util.TreeData
	for _, v := range rows.List {
		tmp := map[string]interface{}{
			"value": v.Data["id"],
			"pid":   v.Data["pid"],
			"label": v.Data["name"],
		}
		tree = append(tree, tmp)
	}
	tree = util.Tree(tree, 0, "pid", "value")

	return tree, nil
}

func (s *service) ExportMenu(ctx *gin.Context, req *menuReq) (*Menu, error) {
	m := s.GetModel(ctx)
	menus := new(Menu)
	row := m.SelectOne(
		db.Field("id", "pid", "module_id", "name", "path", "type", "icon", "page_type", "page_schema", "view", "sort", "status"),
		db.WhereEq("id", req.Id),
	)
	if row.Err != nil {
		return nil, row.Err
	}
	err := row.Binding(menus)
	if err != nil {
		return nil, err
	}
	loopChild(menus, m)

	return menus, nil
}

func (s *service) ImportMenu(ctx *gin.Context, rootid int64, req *Menu) (bool, error) {
	return s.loopCreate(ctx, rootid, req)
}

func (s *service) loopCreate(ctx *gin.Context, rootid int64, menu *Menu) (ok bool, err error) {
	m := s.GetModel(ctx)
	opt := []db.Option{
		db.WhereEq("pid", rootid),
	}

	if rootid != 0 && menu.Type == MenuDIRType {
		return false, errors.New("您正往一个菜单或页面下添加导入一个目录, 此操作是不允许的, 请使用页面中的导入目录功能")
	}

	if rootid == 0 && menu.Type == MenuDIRType {
		opt = append(opt, db.WhereEq("name", menu.Name))
	} else {
		opt = append(opt, db.WhereEq("path", menu.Path))
	}

	exist := m.SelectOne(opt...)
	if exist.Err != nil {
		return false, exist.Err
	}
	update := map[string]interface{}{
		"pid":         rootid,
		"module_id":   menu.ModuleID,
		"name":        menu.Name,
		"path":        menu.Path,
		"type":        menu.Type,
		"icon":        menu.Icon,
		"page_type":   menu.PageType,
		"view":        menu.View,
		"sort":        menu.Sort,
		"status":      menu.Status,
		"page_schema": menu.PageSchema,
	}
	var lastId int64
	if exist != nil {
		id, _ := exist.Get("id")
		_, err = m.Update(update, db.WhereEq("id", id))
		if err != nil {
			return false, err
		}
		lastId = cast.ToInt64(id)
	} else {
		lastId, err = m.Insert(update)
		if err != nil {
			return false, err
		}
	}

	for _, v := range menu.Children {
		_, err = s.loopCreate(ctx, lastId, v)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func loopChild(menu *Menu, m db.Model) {
	if err := menu.GetChildren(m); err == nil {
		for i := range menu.Children {
			loopChild(menu.Children[i], m)
		}
	}
}

func (s *service) MenuMigration(ctx *gin.Context) error {
	m := s.GetModel(ctx)
	list := m.Select(db.WhereEq("type", 2), db.WhereEq("page_type", 0))
	if list.Err != nil {
		return list.Err
	}
	for _, v := range list.List {
		sub := m.Select(db.WhereEq("pid", v.Data["id"]))
		if sub.Err != nil {
			return sub.Err
		}
		path := ""
		schema := util.MapStrInterface{}
		var delId []interface{}
		up := false
		for _, s := range sub.List {
			token := util.String(cast.ToString(s.Data["path"])).Split("/")
			if len(token.Raw()) != 2 { //nolint:gomnd
				continue
			}
			if !util.ArrStr([]string{"form", ":id", "list"}).Has(token.Get(1).Raw()) {
				continue
			}
			up = true
			delId = append(delId, s.Data["id"])
			path = token.Get(0).Raw()
			util.String(util.ToString(s.Data["page_schema"])).DecodeMap()
			schema = schema.Merge(*util.String(util.ToString(s.Data["page_schema"])).DecodeMap())
		}

		if up {
			_, err := m.Update(map[string]interface{}{
				"path":        "/" + path,
				"page_schema": schema.ToString(),
				"page_type":   7,
			}, db.WhereEq("id", v.Data["id"]))
			if err != nil {
				return err
			}
		}
		if len(delId) > 0 {
			_, err := m.Delete(db.WhereIn("id", delId))
			if err != nil {
				return err
			}
		}

	}
	return nil
}
