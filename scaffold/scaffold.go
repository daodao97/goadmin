package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/ecode"
	"github.com/daodao97/goadmin/pkg/sso"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/goadmin/scaffold/dao"
	"github.com/daodao97/xgo/xlog"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/wire"
	"github.com/siddontang/go/num"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"

	"github.com/daodao97/goadmin/pkg/db"
)

var Provider = wire.NewSet(NewUserState, wire.Struct(new(Options), "*"), New)

type Options struct {
	Conf      *Conf      // 脚手架依赖配置
	UserState *UserState // 用户对象
	Cache     cache.Cache
	Sso       *sso.Sso
}

// New 返回脚手架实例
func New(opt *Options) Scaffold {
	_validator := util.NewValidate()
	_validator.RegisterValidation(unique)

	commonConf := NewCommonConfig(opt.Cache)
	schemaPool := NewSchemaPool(commonConf)

	return Scaffold{
		Conf:       opt.Conf,
		Sso:        opt.Sso,
		schemaPool: schemaPool,
		Token: NewToken(&JwtConf{
			Secret:      opt.Conf.Jwt.Secret,
			TokenExpire: opt.Conf.Jwt.TokenExpire,
		}),
		Cache:        opt.Cache,
		validator:    _validator,
		UserState:    opt.UserState,
		CommonConf:   commonConf,
		columnRender: map[string]func(ctx *gin.Context, rows []dao.Row) []dao.Row{},
	}
}

// Scaffold 数据库对象的 crud 脚手架
type Scaffold struct {
	Conf         *Conf
	Sso          *sso.Sso
	Cache        cache.Cache
	TreeList     *TreeList
	Token        *Token
	schemaPool   *SchemaPool
	validator    *util.Validate
	UserState    *UserState
	dao          dao.Dao  // 抽象后的数据层对象, 底层可能为 mysql/es/mongo/api 等等
	model        db.Model // mysql 数据对象, 当 dao==nil 是用此实例
	CommonConf   *CommonConf
	columnRender map[string]func(ctx *gin.Context, rows []dao.Row) []dao.Row
	BeforeCreate func(ctx *gin.Context, val dao.Row) (dao.Row, error)
	BeforeUpdate func(ctx *gin.Context, val dao.Row, id int64) (dao.Row, error)
	AfterCreate  func(ctx *gin.Context, val dao.Row, id int64) error
	AfterUpdate  func(ctx *gin.Context, val dao.Row, id int64) error
	BeforeList   func(ctx *gin.Context, where []dao.Option) []dao.Option
	AfterList    func(ctx *gin.Context, list []dao.Row) []dao.Row
	BeforeGet    func(ctx *gin.Context, val dao.Row) dao.Row
	AfterGet     func(ctx *gin.Context, val dao.Row) dao.Row
	AfterDel     func(ctx *gin.Context, id []int64) error
}

func (s *Scaffold) User(ctx *gin.Context) (*UserAttr, error) {
	token := ctx.Request.Header.Get("x-token")
	if token == "" {
		return nil, fmt.Errorf("x-token in Header not found")
	}

	t, err := s.Token.ParseToken(token)
	if err != nil {
		return nil, ecode.Error(401, "token解析失败")
	}

	m := NewUser()

	row := m.SelectOne(
		db.Field("id", "avatar", "email", "name", "nickname", "role_ids"),
		db.WhereEq("id", t.UserID),
	)
	if row.Err != nil {
		return nil, row.Err
	}

	var info = new(UserAttr)
	err = row.Binding(info)
	if err != nil {
		return nil, err
	}
	info.Env = os.Getenv("DEPLOY_ENV")

	website, _ := s.Website(ctx)
	if website != nil {
		website.Modules = []Module{}
		website.MacroVar = ""
		info.Website = website
	}

	return info, nil
}

// RequestQueryParams 当前请求的 url 参数
func (s *Scaffold) RequestQueryParams(ctx *gin.Context) util.MapStrInterface {
	_query := ctx.Request.URL.Query()
	query := make(map[string]interface{})
	for k, v := range _query {
		parts := strings.Split(v[0], ",")
		if len(parts) > 1 {
			query[k] = parts
			continue
		}

		_k := util.String(k)

		if _k.EndWith("[]") {
			query[_k.ReplaceAll("[]", "").Raw()] = v
		} else {
			query[k] = v[0]
		}
	}
	return query
}

// RequestBody 当前请求的提交数据 支持 application/json, multipart/form-data
func (s *Scaffold) RequestBody(ctx *gin.Context) (util.MapStrInterface, error) {
	input := make(map[string]interface{})
	ctype := ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ctype, binding.MIMEJSON) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil, err
		}
		rdr := ioutil.NopCloser(bytes.NewBuffer(body))
		ctx.Request.Body = rdr
		_ = util.Binding(body, &input)
	} else {
		for k, v := range ctx.Request.Form {
			input[k] = v[0]
		}
	}
	return input, nil
}

func (s *Scaffold) SetDao(d dao.Dao) {
	s.dao = d
}

func (s *Scaffold) Dao(ctx *gin.Context) (dao.Dao, error) {
	if s.dao != nil {
		return s.dao, nil
	}
	m := s.GetModel(ctx)
	return dao.NewMysqlDao(m), nil
}

func (s *Scaffold) SetModel(m db.Model) {
	s.model = m
}

func (s *Scaffold) GetModel(ctx *gin.Context) db.Model {
	return s.model
}

func (s *Scaffold) orderByOpts(ctx *gin.Context, sch *Schema) (opts []dao.Option) {
	if sch.OrderBy != nil {
		opts = append(opts, dao.Orderby(sch.OrderBy.Field, sch.OrderBy.Mod))
	}
	if len(sch.OrderByMulti) > 0 {
		for _, v := range sch.OrderByMulti {
			opts = append(opts, dao.Orderby(v.Field, v.Mod))
		}
	}
	if orderBy := s.orderByOption(ctx); orderBy != nil {
		opts = append(opts, orderBy)
	}

	return opts
}

// List 列表接口
func (s *Scaffold) List(ctx *gin.Context) (*ListResp, error) { //nolint:gocognit
	m, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}
	sch, err := s.getSchema(ctx)
	if err != nil {
		return nil, err
	}

	var opts []dao.Option
	for _, v := range sch.BaseWhere {
		opts = append(opts, dao.Where(v.Field, v.Operator, v.Value))
	}

	// 查询条件构造
	_opts, err := s.listFilter(ctx, sch)
	if err != nil {
		return nil, err
	}

	opts = append(opts, _opts...)

	if s.BeforeList != nil {
		opts = s.BeforeList(ctx, opts)
	}

	// 查询总数
	count, err := m.Count(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// 分页构造
	query := s.RequestQueryParams(ctx)
	ps, exist := query["_ps"]
	if !exist {
		ps = "20"
	}
	pn, exist := query["_pn"]
	if !exist {
		pn = "1"
	}

	if count == 0 {
		return &ListResp{
			Page: Page{
				Ps:    cast.ToInt64(ps),
				Pn:    cast.ToInt64(pn),
				Total: count,
			},
			List: []dao.Row{},
		}, nil
	}

	opts = append(opts, s.orderByOpts(ctx, sch)...)

	var hasMany []Header
	// 查询字段
	var fields []string
	for _, v := range sch.Headers {
		if v.HasMany != "" {
			hasMany = append(hasMany, v)
			continue
		}
		if v.Render != "" || v.Fake {
			continue
		}
		fields = append(fields, v.Field)
	}
	opts = append(opts, dao.Field(fields...))
	opts = append(opts, dao.Pagination(cast.ToInt(pn), cast.ToInt(ps))...)

	// 查询列表
	rows, err := m.Select(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if s.TreeList != nil {
		if !s.TreeList.Lazy {
			var tree []util.TreeData
			for _, v := range rows {
				tree = append(tree, v)
			}
			tree = util.Tree(tree, 0)
			var _rows []dao.Row
			for _, v := range tree {
				_rows = append(_rows, v)
			}
			rows = _rows
		} else {
			for i := range rows {
				rows[i]["hasChildren"] = true
			}
		}
	}

	for _, v := range sch.Headers {
		if v.Render == "" {
			continue
		}
		handle, ok := s.columnRender[v.Render]
		if ok {
			rows = handle(ctx, rows)
		}
	}

	for _, str := range sch.HasOne {
		newRows, err := HasOneData(str, rows)
		if err == nil {
			rows = newRows
		} else {
			xlog.Warn("hasOne", xlog.Err(err), xlog.String("str", str))
		}
	}

	for _, v := range hasMany {
		newRows, err := HasManyData(v.HasMany, v.Field, rows)
		if err == nil {
			rows = newRows
		} else {
			xlog.Warn("hasMany", xlog.Err(err), xlog.String("str", v.HasMany))
		}
	}

	if s.AfterList != nil {
		rows = s.AfterList(ctx, rows)
	}

	resp := &ListResp{
		Page: Page{
			Ps:    cast.ToInt64(ps),
			Pn:    cast.ToInt64(pn),
			Total: count,
		},
		List: rows,
	}

	return resp, nil
}

func (s *Scaffold) listFilter(ctx *gin.Context, sch *Schema) ([]dao.Option, error) {
	var opts []dao.Option
	query := s.RequestQueryParams(ctx)
	f := sch.Filter
	tabFilter := make(map[string]Filter)

	for _, v := range sch.Tabs {
		tabFilter[v.Field] = Filter{
			Label:    v.Label,
			Field:    v.Field,
			Operator: v.Operator,
			Form:     true,
		}
	}

	for _, v := range tabFilter {
		f = append(f, v)
	}

	for i, v := range f {
		if !v.Form {
			continue
		}
		if v.Field == "" {
			return nil, fmt.Errorf("page_schema.tabs.%d 节点配置错误, field 字段是必须的", i)
		}
		val, exist := query[v.Field]
		if !exist {
			continue
		}
		op := "="
		if v.Operator != "" {
			op = v.Operator
		}
		switch val := val.(type) {
		case string:
			switch op {
			case "like":
				opts = append(opts, dao.Where(v.Field, "like", fmt.Sprintf("%%%v%%", val)))
			case "suffix_like":
				opts = append(opts, dao.Where(v.Field, "like", fmt.Sprintf("%v%%", val)))
			case "prefix_like":
				opts = append(opts, dao.Where(v.Field, "like", fmt.Sprintf("%%%v", val)))
			default:
				opts = append(opts, dao.Where(v.Field, op, val))
			}
		case []string:
			switch op {
			case "between":
				opts = append(opts, dao.Where(v.Field, ">=", val[0]))
				opts = append(opts, dao.Where(v.Field, "<=", val[1]))
			default:
				opts = append(opts, dao.Where(v.Field, "in", util.ArrStr(val).ToSliceInterface()))
			}
		}
	}
	if s.TreeList != nil && s.TreeList.Lazy {
		pidKey := "pid"
		if s.TreeList.PidKey != "" {
			pidKey = s.TreeList.PidKey
		}
		opts = append(opts, dao.Where(pidKey, "=", 0))
	}
	return opts, nil
}

func (s *Scaffold) orderByOption(ctx *gin.Context) dao.Option {
	query := s.RequestQueryParams(ctx)
	if orderBy, ok := query["_sort_by"]; ok && cast.ToString(orderBy) != "" {
		mod := "desc"
		_orderMod, ok := query["_sort_type"]
		orderMod := cast.ToString(_orderMod)
		if ok && orderMod != "" && (orderMod == "desc" || orderMod == "asc") {
			mod = orderMod
		}
		return dao.Orderby(cast.ToString(orderBy), mod)
	}
	return nil
}

// Children 树形列表的子项接口
func (s *Scaffold) Children(ctx *gin.Context) ([]dao.Row, error) {
	pid, ok := s.RequestQueryParams(ctx)["pid"]
	if !ok {
		return nil, fmt.Errorf("拉取子项时必须指定父级 pid=xx")
	}
	sch, err := s.getSchema(ctx)
	if err != nil {
		return nil, err
	}
	var fields []string
	for _, v := range sch.Headers {
		fields = append(fields, v.Field)
	}
	opts := []dao.Option{
		dao.Where("pid", "=", pid),
		dao.Field(fields...),
	}

	opts = append(opts, s.orderByOpts(ctx, sch)...)

	m, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := m.Select(ctx, opts...)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i]["hasChildren"] = true
	}
	return rows, nil
}

func (s *Scaffold) Tree(ctx *gin.Context, opt ...dao.Option) ([]util.TreeData, error) {
	m, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := m.Select(ctx, opt...)
	if err != nil {
		return nil, err
	}
	var tree []util.TreeData
	for _, v := range rows {
		tree = append(tree, v)
	}
	return util.Tree(tree, 0), nil
}

// Get 单条记录获取接口
func (s *Scaffold) Get(ctx *gin.Context) (dao.Row, error) {
	id, exist := ctx.Params.Get("id")
	if !exist {
		return nil, errors.New("获取记录必须有主键 如 /user/1")
	}
	row, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.AfterGet != nil {
		row = s.AfterGet(ctx, row)
	}
	return row, nil
}

func (s *Scaffold) GetById(ctx *gin.Context, id string) (dao.Row, error) {
	m, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}

	opts := []dao.Option{
		dao.Where("id", "=", id),
	}

	sch, err := s.getSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("未发现页面配置, 或页面未启用, id:%s", id)
	}

	var fields []string
	for _, v := range sch.FormItems {
		if !v.Form {
			continue
		}
		fields = append(fields, v.Field)
	}
	opts = append(opts, dao.Field(fields...))

	row, err := m.SelectOne(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if s.AfterGet != nil {
		row = s.AfterGet(ctx, row)
	}
	return row, nil
}

// Set 创建 或 更新 接口
func (s *Scaffold) Set(ctx *gin.Context) (int64, error) {
	schema, err := s.getSchema(ctx)
	if err != nil {
		return 0, err
	}
	if schema.Maintenance {
		return 0, errors.New("当前功能正在维护, 请稍后再试")
	}
	_, exist := ctx.Params.Get("id")
	if exist {
		return s.update(ctx)
	}
	return s.insert(ctx)
}

// insert 创建记录
func (s *Scaffold) insert(ctx *gin.Context) (int64, error) {
	m, err := s.Dao(ctx)
	if err != nil {
		return 0, err
	}

	value, err := s.InputData(ctx)
	if err != nil {
		return 0, err
	}
	if s.BeforeCreate != nil {
		value, err = s.BeforeCreate(ctx, value)
		if err != nil {
			return 0, err
		}
	}

	id, err := m.Insert(ctx, value)
	if err != nil {
		return 0, err
	}

	if s.AfterCreate != nil {
		err = s.AfterCreate(ctx, value, cast.ToInt64(id))
		if err != nil {
			return 0, err
		}
	}
	// 记录日志
	// s.AddLogs(ctx, []string{cast.ToString(id)}, "add", nil, value, nil)
	return id, nil
}

// update 更新记录
func (s *Scaffold) update(ctx *gin.Context) (int64, error) {
	_, err := s.Get(ctx)
	if err != nil {
		return 0, err
	}

	value, err := s.InputData(ctx)
	if err != nil {
		return 0, err
	}

	id, _ := ctx.Params.Get("id")
	opts := []dao.Option{
		dao.Where("id", "=", id),
	}

	//before, _ := s.GetById(ctx, id)

	id64 := cast.ToInt64(id)

	if s.BeforeUpdate != nil {
		value, err = s.BeforeUpdate(ctx, value, id64)
		if err != nil {
			return 0, err
		}
	}

	m, _ := s.Dao(ctx)
	_, err = m.Update(ctx, value, opts...)
	if err != nil {
		return 0, err
	}

	if s.AfterUpdate != nil {
		err = s.AfterUpdate(ctx, value, id64)
		if err != nil {
			return 0, err
		}
	}

	// 记录日志
	// s.AddLogs(ctx, []string{id}, "modify", before, value, nil)

	return id64, nil
}

// InputData 经过校验后的用户提交数据
func (s *Scaffold) InputData(ctx *gin.Context) (dao.Row, error) {
	input, err := s.RequestBody(ctx)
	if err != nil {
		return nil, err
	}
	if len(input) == 0 {
		return nil, fmt.Errorf("未获取到提交参数")
	}

	sch, err := s.getSchema(ctx)
	if err != nil {
		return nil, errors.New("not found page_schema")
	}

	value := make(dao.Row)
	for _, v := range sch.FormItems {
		if !v.Form {
			continue
		}
		val, exist := input[v.Field]
		if v.Validate != "" {
			err = s.validator.VarCtx(ctx, val, v.Validate, v.Label)
			if err != nil {
				return nil, err
			}
		}
		if !exist {
			continue
		}
		strVal := cast.ToString(val)
		if v.Type == "json" && strVal != "" {
			strVal, err = util.JsonStrRemoveComments(strVal)
			if err != nil {
				return nil, fmt.Errorf("%s json格式错误, 请检查修正", v.Field)
			}
			isJson := gjson.Valid(strVal)
			if !isJson {
				return nil, fmt.Errorf("%s json格式错误, 请检查修正", v.Field)
			}
		}
		value[v.Field] = val
	}

	if len(value) == 0 {
		return nil, fmt.Errorf("未解析到数据")
	}

	return value, nil
}

// Del 删除记录
func (s *Scaffold) Del(ctx *gin.Context) (bool, error) {
	schema, err := s.getSchema(ctx)
	if err != nil {
		return false, err
	}
	if schema.Maintenance {
		return false, errors.New("当前功能正在维护, 请稍后再试")
	}
	m := s.GetModel(ctx)
	id, exist := s.RequestQueryParams(ctx)["id"]
	if !exist {
		return false, errors.New("必须指定删除记录id")
	}
	ids := strings.Split(cast.ToString(id), ",")
	var _ids []interface{}
	for _, v := range ids {
		_ids = append(_ids, v)
	}
	opts := []db.Option{db.WhereIn("id", _ids)}
	_, err = m.Delete(opts...)
	if err != nil {
		return false, err
	}
	if s.AfterDel != nil {
		var tmp []int64
		for _, i := range _ids {
			tmp = append(tmp, cast.ToInt64(i))
		}
		err = s.AfterDel(ctx, tmp)
		if err != nil {
			return false, err
		}
	}
	// 记录日志
	// s.AddLogs(ctx, ids, "del", nil, nil, nil)
	return true, nil
}

func (s *Scaffold) currentPath(ctx *gin.Context) string {
	return ctx.Request.Header.Get("x-path")
}

// getSchema 对外输出当前路由自定解析的 path 结构化的 Schema 数据
func (s *Scaffold) getSchema(ctx *gin.Context) (*Schema, error) {
	project := s.currentPath(ctx)
	sch, err := s.schemaPool.GetSchema(project)
	if err != nil {
		return nil, fmt.Errorf("未找到响应页面的PageSchema 或 页面未启用 id:%s  err:%v", project, err)
	}
	return sch, nil
}

// GetSchemaByRoute 对外输出指定path全量的 page_schema 数据
func (s *Scaffold) GetSchemaByRoute(ctx *gin.Context, route string) (interface{}, error) {
	sch, ok := s.schemaPool.Get(route)
	if !ok {
		return nil, fmt.Errorf("未找到响应页面的PageSchema 或 页面未启用 id:%s", route)
	}
	return sch, nil
}

func (s *Scaffold) Website(ctx *gin.Context) (*Website, error) {
	var website = new(Website)
	err := s.CommonConf.Get(ctx, "website", website)
	if err != nil {
		return nil, err
	}
	return website, nil
}

// SelectOption 下拉框备选项的通用搜索
func (s *Scaffold) SelectOption(ctx *gin.Context, kw string, valueKey string, labelKey string, opts ...dao.Option) ([]dao.Row, error) {
	m, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}
	return s.SelectOptionWithModel(ctx, m, kw, valueKey, labelKey, opts...)
}

// SelectOptionWithModel 下拉框备选项的通用搜索
func (s *Scaffold) SelectOptionWithModel(ctx *gin.Context, m dao.Dao, kw string, valueKey string, labelKey string, opts ...dao.Option) ([]dao.Row, error) {
	strs := strings.Split(kw, ",")
	var _ids []interface{}
	for _, str := range strs {
		id, err := num.ParseInt64(str)
		if err != nil {
			break
		}
		_ids = append(_ids, id)
	}
	if len(_ids) == len(strs) {
		if valueKey == "" {
			valueKey = "id"
		}
		if len(_ids) == 1 {
			opts = append(opts, dao.Where(valueKey, "=", _ids[0]))
		} else {
			opts = append(opts, dao.Where(valueKey, "in", _ids))
		}
	} else {
		if labelKey == "" {
			labelKey = "title"
		}
		opts = append(opts, dao.Where(labelKey, "like", fmt.Sprintf("%%%s%%", kw)))
	}
	opts = append(opts, dao.Limit(50))

	list, err := m.Select(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Scaffold) RegColumnRender(name string, fn func(ctx *gin.Context, rows []dao.Row) []dao.Row) {
	s.columnRender[name] = fn
}
