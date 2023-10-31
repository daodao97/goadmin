package scaffold

import (
	"fmt"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/goadmin/scaffold/dao"
	"regexp"
	"strings"

	"github.com/spf13/cast"

	"github.com/daodao97/goadmin/pkg/db"
)

type hasOpts struct {
	Pool       string
	DB         string
	Table      string
	LocalKey   string
	ForeignKey string
	OtherKeys  []string
}

// [pool.]db.table:[local_key->]foreign_key,other_key
func explodeHasStr(str string) (opt *hasOpts, err error) {
	var re = regexp.MustCompile(`([a-zA-Z_0-9]+\.)?([a-zA-Z_0-9]+)\.([a-zA-Z_0-9]+):([a-zA-Z_0-9]+->)?([a-zA-Z_0-9 ]+)?([a-zA-Z_,0-9 ]+)`)
	if !re.MatchString(str) {
		return nil, fmt.Errorf("has string syntux is error, mast like [pool.]db.table:[local_key->]foreign_key,other_key")
	}
	matched := re.FindStringSubmatch(str)
	for i, v := range matched {
		if i == 0 {
			continue
		}
		matched[i] = strings.ReplaceAll(strings.ReplaceAll(v, "->", ""), ".", "")
	}

	opt = &hasOpts{
		Pool:       matched[1],
		DB:         matched[2],
		Table:      matched[3],
		LocalKey:   matched[4],
		ForeignKey: matched[5],
		OtherKeys: util.String(matched[6]).Split(",").Filter(func(index int, str string) bool {
			return str != ""
		}).Raw(),
	}
	if opt.Pool == "" {
		opt.Pool = "default"
	}
	if opt.LocalKey == "" {
		opt.LocalKey = "id"
	}
	if opt.ForeignKey == "" {
		opt.ForeignKey = "id"
	}

	return opt, nil
}

func HasOneData(hasStr string, list []dao.Row) (newList []dao.Row, err error) {
	if len(list) == 0 {
		return
	}
	opt, err := explodeHasStr(hasStr)
	if err != nil {
		return nil, err
	}
	model := db.New(opt.Table, db.WithConn(opt.Pool))
	fields := append(opt.OtherKeys, opt.ForeignKey)
	var ids []interface{}
	for _, v := range list {
		ids = append(ids, v[opt.LocalKey])
	}
	res := model.Select(db.WhereIn(opt.ForeignKey, ids), db.Field(fields...))
	if res.Err != nil {
		return nil, res.Err
	}

	_list := make(map[string]db.Row)
	for _, o := range res.List {
		_list[cast.ToString(o.Data[opt.ForeignKey])] = o
	}

	otherKeys := make([]string, len(opt.OtherKeys))
	for i, k := range opt.OtherKeys {
		_k := util.String(k).RegexSplit(`(?i) +as +`).Map(func(str string) string {
			return util.String(str).TrimSpace().Raw()
		}).Last().Raw()
		otherKeys[i] = _k
	}

	for i, l := range list {
		if o, ok := _list[cast.ToString(l[opt.LocalKey])]; ok {
			for _, k := range otherKeys {
				list[i][k] = o.Data[k]
			}
		}
	}

	return list, nil
}

func HasManyData(hasStr string, key string, list []dao.Row) (newList []dao.Row, err error) {
	if len(list) == 0 {
		return
	}
	opt, err := explodeHasStr(hasStr)
	if err != nil {
		return nil, err
	}
	model := db.New(opt.Table, db.WithConn(opt.Pool))
	if err != nil {
		return nil, err
	}
	fields := append(opt.OtherKeys, opt.ForeignKey)
	var ids []interface{}
	for _, v := range list {
		ids = append(ids, v[opt.LocalKey])
	}
	res := model.Select(db.WhereIn(opt.ForeignKey, ids), db.Field(fields...))
	if res.Err != nil {
		return nil, res.Err
	}

	otherKeys := make([]string, len(opt.OtherKeys))
	for i, k := range opt.OtherKeys {
		_k := util.String(k).RegexSplit(`(?i) +as +`).Map(func(str string) string {
			return util.String(str).TrimSpace().Raw()
		}).Last().Raw()
		otherKeys[i] = _k
	}

	_filter := func(keys []string, data db.Row) db.Row {
		newData := db.Row{Data: map[string]interface{}{}}
		for k, v := range data.Data {
			if util.ArrStr(keys).Has(k) {
				newData.Data[k] = v
			}
		}
		return newData
	}

	_list := make(map[string][]db.Row)
	for _, o := range res.List {
		_list[cast.ToString(o.Data[opt.ForeignKey])] = append(_list[cast.ToString(o.Data[opt.ForeignKey])], _filter(otherKeys, o))
	}

	for i, l := range list {
		if o, ok := _list[cast.ToString(l[opt.LocalKey])]; ok {
			list[i][key] = o
		}
	}

	return list, nil
}
