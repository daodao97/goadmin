package scaffold

import (
	"fmt"
	"strings"
	"time"

	"github.com/daodao97/goadmin/pkg/util"
)

func filterByEnv(data, envIdent string) string { //nolint:gocognit
	index := strings.Index(data, envIdent)
	if index == -1 {
		return data
	}

	startIndex := index
	preloop := 0
	for {
		startIndex--
		if string(data[startIndex]) == "}" {
			preloop++
			continue
		}
		if preloop > 0 && string(data[startIndex]) == "{" {
			preloop--
			continue
		}
		if string(data[startIndex]) == "{" {
			break
		}

	}
	endIndex := index + len(envIdent)
	inloop := 0
	for {
		endIndex++
		if string(data[endIndex]) == "{" {
			inloop++
			continue
		}
		if inloop > 0 && string(data[endIndex]) == "}" {
			inloop--
			continue
		}
		if string(data[endIndex]) == "}" {
			break
		}

	}

	if startIndex == 0 {
		return "{}"
	}
	preComma := string(data[startIndex-1])
	endComma := string(data[endIndex+1])
	if preComma == "," && endComma == "," {
		startIndex--
	}
	if preComma != "," && endComma == "," {
		endIndex++
	}
	if preComma == "," && endComma != "," {
		startIndex--
	}

	str := data[:startIndex]
	strend := data[endIndex+1:]

	return filterByEnv(str+strend, envIdent)
}

// NewSchemaPool 页面定义的对象池
func NewSchemaPool(cc *CommonConf) *SchemaPool {
	return &SchemaPool{
		pool:       util.NewPool(),
		commonConf: cc,
	}
}

type cacheSchema struct {
	schema   interface{}
	expireAt time.Time
}

type SchemaPool struct {
	pool       *util.Pool
	commonConf *CommonConf
}

// GetSchema 返回 结构体 供 scaffold list/update 等接口使用
func (p *SchemaPool) GetSchema(key string) (*Schema, error) {
	schema, err := p.getSchema(key)
	if err != nil {
		return nil, err
	}
	_s := new(Schema)
	err = util.Binding(schema, _s)
	if err != nil {
		return nil, err
	}

	return _s, nil
}

// getSchema 返回 全量页面配置的json数据, 透传给前端
func (p *SchemaPool) getSchema(route string) (interface{}, error) {
	return util.NewPipErr().
		Wrap(func(input interface{}) (interface{}, error) {
			return SchemaMacroVal(p.commonConf)
		}).
		Wrap(func(input interface{}) (interface{}, error) {
			return getSchemaByRoute(route, *(input).(*map[string]interface{}))
		}).
		Wrap(func(input interface{}) (interface{}, error) {
			jsonStr := input.(string)
			envs := func() []string {
				if util.IsProd() {
					return []string{"pre", "uat"}
				}
				if util.IsPre() {
					return []string{"prod", "uat"}
				}
				if util.IsUat() {
					return []string{"pre", "prod"}
				}
				return []string{}
			}()

			for i := range envs {
				jsonStr = filterByEnv(jsonStr, fmt.Sprintf(`"_env":"%s"`, envs[i]))
			}

			return jsonStr, nil
		}).
		Wrap(func(input interface{}) (interface{}, error) {
			s := new(interface{})
			err := util.Binding(input, s)
			return s, err
		}).
		Run()
}

func (p *SchemaPool) Set(key string, db *Schema) {
	p.pool.Set(key, db)
}

func (p *SchemaPool) Get(key string) (interface{}, bool) {
	v, ok := p.pool.Get(key)
	if ok && v.(cacheSchema).expireAt.Unix() > time.Now().Unix() {
		return v.(cacheSchema).schema, true
	}
	s, err := p.getSchema(key)
	if err != nil {
		return nil, false
	}
	p.pool.Set(key, cacheSchema{
		schema:   s,
		expireAt: time.Now().Add(10 * time.Second),
	})
	return s, true
}
