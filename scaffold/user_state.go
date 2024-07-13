package scaffold

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/util"

	"github.com/daodao97/goadmin/pkg/db"
)

const tokenKey = "admin:user_login_tokens:%d"

const RoleTypeFull = 0
const RoleTypeReadonly = 1

func NewUserState(_cache cache.Cache, conf *Conf) (*UserState, error) {
	return &UserState{
		client: _cache,
		conf:   conf,
		user:   NewUser(),
		role:   db.New("role", db.ColumnHook(db.Json("config"), db.Json("resource"))),
	}, nil
}

type UserState struct {
	client cache.Cache
	conf   *Conf
	user   db.Model
	role   db.Model
}

func (s UserState) HavePermission(ctx context.Context, id int, req *http.Request) bool {
	if supper, _ := s.IsSupper(id); supper {
		return true
	}
	conf, err := s.userConfigs(id)
	if err != nil || len(conf) == 0 {
		return false
	}
	var handle = func(r role) bool {
		// todo 判断当前path是否为该角色所拥有资源
		// 判断是否只读
		if !util.ArrStr([]string{http.MethodGet, http.MethodOptions}).Has(req.Method) &&
			r.Config != nil &&
			util.ArrInt64(r.Config.RoleType).Has(RoleTypeReadonly) {
			return false
		}
		return true
	}
	tmp := 0
	for _, v := range conf {
		if handle(v) {
			tmp += 1
		}
	}

	return tmp > 0
}

func (s UserState) userRoleIds(id int) []int {
	row := s.user.SelectOne(db.WhereEq("id", id))
	var _roleIds []int
	if row.Err != nil || row.Data == nil {
		return _roleIds
	}
	roleIds := new([][]int)
	err := util.Binding(row.Data["role_ids"], roleIds)
	if err != nil {
		return _roleIds
	}
	for _, v := range *roleIds {
		if len(v) > 0 {
			_roleIds = append(_roleIds, v[len(v)-1])
		}
	}
	return _roleIds
}

func (s UserState) userConfigs(id int) ([]role, error) {
	roleIds := s.userRoleIds(id)
	if len(roleIds) == 0 {
		return []role{}, nil
	}
	_rows := new([]role)
	rows := s.role.Select(db.WhereIn("id", arrIntToInterface(roleIds)))
	err := rows.Binding(_rows)
	if err != nil {
		return nil, err
	}
	return *_rows, nil
}

func (s UserState) IsSupper(id int) (bool, error) {
	row := s.user.SelectOne(
		db.Field("id", "avatar", "email", "name", "nickname", "role_ids"),
		db.WhereEq("id", id),
	)
	if row.Err != nil {
		return false, row.Err
	}

	var info = new(UserAttr)
	err := row.Binding(info)
	if err != nil {
		return false, err
	}
	return info.IsSupper(), nil
}

func (s UserState) IsLogin(ctx context.Context, id int) (bool, error) {
	res, err := s.client.Get(ctx, fmt.Sprintf(tokenKey, id))
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return false, err
	}
	return res != "", nil
}

func (s UserState) SetToken(ctx context.Context, id int, token string) error {
	return s.client.SetWithTTL(ctx, fmt.Sprintf(tokenKey, id), token, time.Duration(s.conf.Jwt.TokenExpire)*time.Second)
}

func (s UserState) DelToken(ctx context.Context, id int) error {
	return s.client.Del(ctx, fmt.Sprintf(tokenKey, id))
}

func (s UserState) SetCaptcha(ctx context.Context, captcha string) error {
	return s.client.SetWithTTL(ctx, s.captchaKey(captcha), "1", time.Duration(120)*time.Second)
}

func (s UserState) captchaKey(captcha string) string {
	return fmt.Sprintf("oms:captcha:%s", captcha)
}

func (s UserState) ExistCaptcha(ctx context.Context, captcha string) (bool, error) {
	res, err := s.client.Get(ctx, s.captchaKey(captcha))
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return false, err
	}
	return res != "", nil
}

func arrIntToInterface(arr []int) []interface{} {
	var tmp []interface{}
	for _, v := range arr {
		tmp = append(tmp, v)
	}
	return tmp
}
