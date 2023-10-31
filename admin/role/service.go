package role

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"

	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
	"github.com/gin-gonic/gin"
)

func newService(s *scaffold.Scaffold) *service {
	s.SetModel(db.New("role", db.ColumnHook(db.Json("resource"), db.Json("config"))))
	s.TreeList = &scaffold.TreeList{
		PidKey: "pid",
	}
	return &service{
		Scaffold: s,
	}
}

type service struct {
	*scaffold.Scaffold
}

func (e *service) GetScaffold() *scaffold.Scaffold {
	return e.Scaffold
}

func (e *service) Tree(ctx *gin.Context) ([]util.TreeData, error) {
	m := e.GetModel(ctx)
	rows := m.Select(db.Field("id", "pid", "name"), db.WhereEq("status", 1))
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

func (e *service) Resource(ctx *gin.Context) (interface{}, error) {
	website, err := e.Website(ctx)
	if err != nil {
		return nil, err
	}
	rows := scaffold.NewPage().Select(db.Field("id", "pid", "module_id", "name", "page_schema"))
	if rows.Err != nil {
		return nil, rows.Err
	}

	var list []util.TreeData
	for _, v := range rows.List {
		tmp := map[string]interface{}{
			"pid":       v.Data["pid"],
			"value":     cast.ToString(v.Data["id"]),
			"label":     v.Data["name"],
			"module_id": v.Data["module_id"],
		}

		//if v["page_schema"] != nil {
		//	schema := v["page_schema"].(*map[string]interface{})
		//	s, err := scaffold.SchemaVarReplace(cast.ToInt(v["id"]), schema)
		//	if err != nil {
		//		continue
		//	}
		//	tmp["children"] = s
		//}

		list = append(list, tmp)
	}

	tree := util.Tree(list, 0, "pid", "value")

	var resource []scaffold.TreeOptions

	for _, v := range website.Modules {
		module := scaffold.TreeOptions{
			Value: fmt.Sprintf("module.%s", cast.ToString(v.Id)),
			Label: v.Label,
		}
		for _, r := range tree {
			if cast.ToInt(v.Id) == cast.ToInt(r["module_id"]) {
				tmp := new(scaffold.TreeOptions)
				err = util.Binding(r, tmp)
				if err != nil {
					return nil, err
				}
				module.Children = append(module.Children, *tmp)
			}
		}
		resource = append(resource, module)
	}
	return resource, nil
}
