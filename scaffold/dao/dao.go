package dao

import (
	"context"
	"encoding/json"
)

type Options struct {
	Table   string
	Field   []string
	Where   []where
	Orderby []orderby
	Groupby string
	Limit   int
	Offset  int
	Value   []interface{}
}

type where struct {
	Field    string
	Operator string
	Value    interface{}
	Logic    string
	Sub      []where
}

type Option interface {
	Apply(opts *Options)
}

type Row map[string]interface{}

func (r *Row) Binding(b interface{}) error {
	bt, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, b)
}

type table string

func (t *table) Apply(opts *Options) {
	opts.Table = string(*t)
}

func Table(name string) Option {
	t := table(name)
	return &t
}

type fields []string

func (o fields) Apply(opts *Options) {
	opts.Field = o
}

func Field(name ...string) Option {
	fs := fields(name)
	return &fs
}

type orderby struct {
	field string
	mod   string
}

func (t *orderby) Apply(opts *Options) {
	exist := false
	for i, v := range opts.Orderby {
		if v.field == t.field {
			exist = true
			opts.Orderby[i].mod = v.mod
		}
	}
	if !exist {
		opts.Orderby = append(opts.Orderby, *t)
	}
}

func Orderby(field string, mod string) Option {
	t := orderby{
		field: field,
		mod:   mod,
	}
	return &t
}

type Offset int

func (o Offset) Apply(opts *Options) {
	opts.Offset = int(o)
}

type Limit int

func (o Limit) Apply(opts *Options) {
	opts.Limit = int(o)
}

func Pagination(pageNumber, pageSize int) []Option {
	return []Option{
		Limit(pageSize),
		Offset((pageNumber - 1) * pageSize),
	}
}

func (o *where) Apply(opts *Options) {
	opts.Where = append(opts.Where, where{
		Field:    o.Field,
		Operator: o.Operator,
		Value:    o.Value,
		Logic:    o.Logic,
		Sub:      o.Sub,
	})
}

func Where(field, operator string, value interface{}) Option {
	return &where{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

type Dao interface {
	Count(ctx context.Context, opt ...Option) (int64, error)
	Select(ctx context.Context, opt ...Option) (rows []Row, err error)
	SelectOne(ctx context.Context, opt ...Option) (row Row, err error)
	Insert(ctx context.Context, record Row) (res int64, err error)
	Update(ctx context.Context, record Row, opt ...Option) (res bool, err error)
	Delete(ctx context.Context, opt ...Option) (res bool, err error)
}
