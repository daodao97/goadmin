package scaffold

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/pkg/cache"
	"github.com/daodao97/goadmin/pkg/util"

	"github.com/daodao97/goadmin/pkg/db"
)

func NewCommonConfig(_cache cache.Cache) *CommonConf {
	return &CommonConf{
		client: _cache,
	}
}

type CommonConf struct {
	client cache.Cache
}

func (c CommonConf) key(name string) string {
	return fmt.Sprintf("admin:common_config:%s", name)
}

func (c CommonConf) get(name string) (string, error) {
	res := NewCommonConfModel().SelectOne(db.Field("value"), db.WhereEq("name", name))
	if res.Err != nil {
		return "", res.Err
	}

	return cast.ToString(res.Data["value"]), nil
}

func (c CommonConf) Get(ctx context.Context, name string, bind interface{}) error {
	bindTo := func(val string) error {
		return util.Binding(val, bind)
	}

	val, _ := c.client.Get(ctx, c.key(name))
	if val != "" {
		return bindTo(val)
	}

	val, err := c.get(name)
	if err != nil {
		return err
	}

	_ = c.client.Set(ctx, c.key(name), val)

	return bindTo(val)
}

func (c CommonConf) Update(ctx context.Context, name string, value map[string]interface{}) error {
	str, _ := json.Marshal(value)
	return c.client.Set(ctx, c.key(name), string(str))
}

func (c CommonConf) Delete(ctx context.Context, name string) error {
	return c.client.Del(ctx, c.key(name))
}
