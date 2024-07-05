package user

import (
	"errors"
	"fmt"
	"github.com/daodao97/goadmin/pkg/sso"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
	"github.com/daodao97/goadmin/scaffold/dao"
	"github.com/gin-gonic/gin"
)

func newService(s *scaffold.Scaffold) *service {
	s.SetModel(db.New("user", db.ColumnHook(db.Json("role_ids"))))
	s.RegColumnRender("is_login", func(ctx *gin.Context, rows []dao.Row) []dao.Row {
		for i, row := range rows {
			id := cast.ToInt(row["id"])
			isLogin, _ := s.UserState.IsLogin(ctx, id)
			rows[i]["is_login"] = 0
			if isLogin {
				rows[i]["is_login"] = 1
			}
		}
		return rows
	})
	s.AfterCreate = func(ctx *gin.Context, val dao.Row, id int64) error {
		return nil
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

func (s *service) Routes(ctx *gin.Context) (interface{}, error) {
	website, err := s.Website(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.User(ctx)
	if err != nil {
		return nil, err
	}

	resourceIds, err := s.resourceIds(user)
	if err != nil {
		return nil, err
	}

	if user.IsSupper() {
		resourceIds = []string{"0"}
	}

	resource, err := s.resource(resourceIds)
	if err != nil {
		return nil, err
	}

	var list []util.TreeData
	for _, v := range resource {
		list = append(list, v.Data)
	}

	tree := util.Tree2(list, util.WithHook(func(data []util.TreeData) []util.TreeData {
		sort.Slice(data, func(i, j int) bool {
			return cast.ToInt(data[i]["sort"]) > cast.ToInt(data[j]["sort"])
		})
		return data
	}))

	for i, v := range website.Modules {
		for _, r := range tree {
			if cast.ToInt(v.Id) == cast.ToInt(r["module_id"]) {
				route := new(scaffold.Route)
				err = util.Binding(r, route)
				if err != nil {
					return nil, fmt.Errorf("Route Binding err %v ", err)
				}
				website.Modules[i].Routes = append(website.Modules[i].Routes, *route)
			}
		}
	}

	var m []scaffold.Module
	for _, v := range website.Modules {
		if len(v.Routes) > 0 {
			m = append(m, v)
		}
	}

	return m, nil
}

func (s *service) resourceIds(user *scaffold.UserAttr) (util.ArrStr, error) {
	var roleIds []interface{}
	for _, v := range user.RoleIds {
		if v == nil {
			continue
		}
		roleIds = append(roleIds, v[len(v)-1])
	}
	if len(roleIds) == 0 {
		return []string{}, nil
	}

	m := db.New("role", db.ColumnHook(db.Json("resource")))

	roles := m.Select(db.Field("resource"), db.WhereIn("id", roleIds))
	if roles.Err != nil {
		return nil, roles.Err
	}

	resource := new([]struct {
		Resource []string `json:"resource"`
	})
	err := roles.Binding(resource)
	if err != nil {
		return nil, err
	}

	_resource := util.ArrStr([]string{})
	for _, v := range *resource {
		_resource = _resource.Concat(v.Resource)
	}

	return _resource, nil
}

func (s *service) parentDirMenuIds(ids util.ArrStr) util.ArrStr {
	opt := []db.Option{
		db.Field("id", "pid"),
		db.WhereEq("status", 1),
		db.WhereIn("type", []interface{}{1, 2}),
		db.WhereIn("id", ids.Unique().ToSliceInterface()),
	}
	parent := scaffold.NewPage().Select(opt...)
	if len(parent.List) == 0 {
		return ids
	}
	var parentIds util.ArrStr
	for _, v := range parent.List {
		parentIds = append(parentIds, cast.ToString(v.Data["pid"]))
	}
	parentIds = parentIds.Filter(func(index int, str string) bool {
		return str != ""
	})
	if parentIds.Length() == 0 {
		return ids
	}
	return ids.Concat(parentIds, s.parentDirMenuIds(parentIds)).Unique()
}

func (s *service) allResourceIds(ids util.ArrStr) util.ArrStr {
	opt := []db.Option{
		db.Field("id", "pid"),
		db.WhereEq("status", 1),
	}
	var moduleIds []interface{}
	ids = ids.Map(func(str string) string {
		_str := util.String(str)
		if !_str.Contains(".") {
			return str
		}
		parts := _str.Split(".")
		if _str.StartWith("module.") {
			moduleIds = append(moduleIds, parts.Last())
			return ""
		}
		return parts.First().Raw()
	}).Unique().Filter(func(index int, str string) bool {
		return str != ""
	})
	var group []db.Option
	if ids.Length() > 0 {
		group = append(group, db.WhereIn("pid", ids.ToSliceInterface()))
	}
	if len(moduleIds) > 0 {
		if len(group) == 0 {
			group = append(group, db.WhereIn("module_id", moduleIds))
		} else {
			group = append(group, db.WhereOrIn("module_id", moduleIds))
		}
	}
	opt = append(opt, db.WhereGroup(group...))
	sub := scaffold.NewPage().Select(opt...)
	if len(sub.List) == 0 {
		return ids
	}
	var subids util.ArrStr
	for _, v := range sub.List {
		subids = append(subids, cast.ToString(v.Data["id"]))
	}
	return ids.Concat(subids, s.allResourceIds(subids)).Unique()
}

func (s *service) resource(ids util.ArrStr) ([]db.Row, error) {
	if ids.Length() == 0 {
		return []db.Row{}, nil
	}

	opt := []db.Option{
		db.Field("id", "pid", "module_id", "name", "icon", "path", "page_type", "type", "view", "sort"),
		db.WhereEq("status", 1),
		//db.Orderby("sort", db.DESC),
	}
	allids := s.allResourceIds(ids).Unique()
	allids = s.parentDirMenuIds(allids.Concat(ids)).Unique()
	if allids.Length() == 0 {
		return []db.Row{}, nil
	}

	opt = append(opt, db.WhereIn("id", allids.ToSliceInterface()))
	rows := scaffold.NewPage().Select(opt...)
	if rows.Err != nil {
		return nil, rows.Err
	}

	return rows.List, nil
}

func (s *service) Login(ctx *gin.Context) (*LoginRes, error) {
	req := new(LoginReq)
	if err := ctx.ShouldBindWith(req, binding.JSON); err != nil {
		return nil, err
	}

	ssoMap := *s.Scaffold.Sso

	var row *db.Row

	// sso 登录
	if req.Ticket != "" && ssoMap != nil {
		_sso, ok := ssoMap[req.Key]
		if !ok {
			return nil, errors.New(fmt.Sprintf("不支持的统一登录渠道 %s", req.Key))
		}

		ssoInfo, err := _sso.GetUserInfo(ctx, req.Ticket)
		if err != nil {
			return nil, err
		}

		_row, err := s.getUserBySsoInfo(ctx, ssoInfo)
		if err != nil {
			return nil, err
		}

		row = _row
		if row.Err != nil && !errors.Is(row.Err, db.ErrNotFound) {
			return nil, row.Err
		}
	} else {
		if exist, err := s.UserState.ExistCaptcha(ctx, strings.ToLower(req.Captcha)); err != nil || !exist {
			return nil, fmt.Errorf("验证码错误或已过期")
		}
		md5 := LoginMD5(req.Username, req.Password, strings.ToLower(req.Captcha))
		if md5 != req.Sing {
			return nil, fmt.Errorf("校验错误")
		}

		opts := []db.Option{
			db.WhereEq("name", req.Username),
		}
		m := s.GetModel(ctx)
		row = m.SelectOne(opts...)
		if row.Err != nil && !errors.Is(row.Err, db.ErrNotFound) {
			return nil, row.Err
		}
		err := bcrypt.CompareHashAndPassword([]byte(cast.ToString(row.Data["password"])), []byte(req.Password))
		if err != nil {
			return nil, fmt.Errorf("密码错误")
		}
	}

	token, err := s.genToken(ctx, cast.ToInt(row.Data["id"]), cast.ToString(row.Data["email"]))
	if err != nil {
		return nil, err
	}
	return &LoginRes{
		Name:  cast.ToString(row.Data["nickname"]),
		Token: token,
	}, nil
}

func (s *service) genToken(ctx *gin.Context, id int, email string) (token string, err error) {
	token, err = s.Token.GenerateToken(id, email)
	if err != nil {
		return "", err
	}
	err = s.UserState.SetToken(ctx, id, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) Logout(ctx *gin.Context) error {
	user, err := s.User(ctx)
	if err != nil {
		return err
	}
	err = s.UserState.DelToken(ctx, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Offline(ctx *gin.Context) error {
	id, ok := ctx.Params.Get("id")
	if !ok {
		return fmt.Errorf("缺少必要参数")
	}
	err := s.UserState.DelToken(ctx, cast.ToInt(id))
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdatePwd(ctx *gin.Context) error {
	id, ok := ctx.Params.Get("id")
	if !ok {
		return fmt.Errorf("缺少必要参数")
	}
	val, err := s.RequestBody(ctx)
	if err != nil {
		return err
	}

	pwd, ok := val["password"]
	if !ok || pwd == "" {
		return fmt.Errorf("缺少必要参数")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cast.ToString(pwd)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	m := s.GetModel(ctx)
	_, err = m.Update(map[string]interface{}{"password": string(hash)}, db.WhereEq("id", id))
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RoutesSE(ctx *gin.Context) (interface{}, error) {
	rTmp, err := s.Routes(ctx)
	if err != nil {
		return nil, err
	}
	rTmp2, ok := rTmp.([]scaffold.Module)
	if !ok {
		return nil, fmt.Errorf("RoutesSE 初始化数据错误")
	}
	var resq []*scaffold.Route
	for _, v := range rTmp2 {
		if len(v.Routes) == 0 {
			continue
		}
		resq = append(resq, s.routerChildren(v.Routes)...)
	}
	return resq, nil
}

func (s *service) routerChildren(rs []scaffold.Route) (res []*scaffold.Route) {
	for _, v := range rs {
		if v.ID == 0 {
			continue
		}
		res = append(res, &scaffold.Route{ID: v.ID, Name: v.Name})
		for _, v2 := range s.routerChildren(v.Children) {
			res = append(res, &scaffold.Route{ID: v2.ID, Name: fmt.Sprintf("%s / %s", v.Name, v2.Name)})
		}
	}
	return
}

func (s *service) getUserBySsoInfo(ctx *gin.Context, info *sso.UserInfo) (row *db.Row, err error) {
	opts := []db.Option{
		db.WhereEq("name", info.Name),
	}
	m := s.GetModel(ctx)
	row = m.SelectOne(opts...)
	if row.Err != nil && !errors.Is(row.Err, db.ErrNotFound) {
		return nil, row.Err
	}
	if row.Data != nil {
		return row, nil
	}
	record := db.Record{
		"name":     info.Name,
		"nickname": info.Nickname,
		"status":   1,
		"email":    info.Email,
	}
	if info.DefaultRoleId > 0 {
		record["role_ids"] = [][]int{{info.DefaultRoleId}}
	} else {
		record["role_ids"] = [][]int{{}}
	}
	id, err := m.Insert(record)
	if err != nil {
		return nil, err
	}
	row = m.SelectOne(db.WhereEq("id", id))
	return row, nil
}
