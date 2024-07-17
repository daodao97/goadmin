package dao

import (
	"context"

	"github.com/daodao97/goadmin/pkg/db"
)

func NewMysqlDao(m db.Model) Dao {
	return &mysqldao{
		model: m,
	}
}

type mysqldao struct {
	model db.Model
}

func (m *mysqldao) transOptions(opts ...Option) (options []db.Option) {
	daoOptions := new(Options)
	for _, v := range opts {
		v.Apply(daoOptions)
	}
	if daoOptions == nil {
		return []db.Option{}
	}

	opt := []db.Option{
		db.Field(daoOptions.Field...),
		db.Limit(daoOptions.Limit),
		db.Offset(daoOptions.Offset),
		db.Value(daoOptions.Value...),
	}

	for _, v := range daoOptions.Orderby {
		if v.field != "" && v.mod != "" {
			if v.mod == "desc" {
				opt = append(opt, db.OrderByDesc(v.field))
			} else {

				opt = append(opt, db.OrderByAsc(v.field))
			}
		}
	}

	// todo 嵌套类型的 Where
	whereOps := func(w ...where) []db.Option {
		var _w []db.Option
		for _, v := range w {
			o := "="
			if v.Operator != "" {
				o = v.Operator
			}
			_w = append(_w, db.Where(v.Field, o, v.Value))
		}
		return _w
	}(daoOptions.Where...)

	if len(whereOps) > 0 {
		opt = append(opt, whereOps...)
	}

	return opt
}

func (m *mysqldao) Count(ctx context.Context, opt ...Option) (int64, error) {
	return m.model.Count(m.transOptions(opt...)...)
}

func (m *mysqldao) Select(ctx context.Context, opt ...Option) ([]Row, error) {
	rows := m.model.Select(m.transOptions(opt...)...)
	if rows.Err != nil {
		return nil, rows.Err
	}
	var list []Row
	for _, v := range rows.List {
		list = append(list, v.Data)
	}
	return list, nil
}

func (m *mysqldao) SelectOne(ctx context.Context, opt ...Option) (Row, error) {
	opt = append(opt, Limit(1), Offset(0))
	row := m.model.SelectOne(m.transOptions(opt...)...)
	if row.Err != nil {
		return nil, row.Err
	}
	return row.Data, nil
}

func (m *mysqldao) Insert(ctx context.Context, record Row) (int64, error) {
	LastInsertId, err := m.model.Insert(db.Record(record))
	if err != nil {
		return 0, err
	}
	return LastInsertId, nil
}

func (m *mysqldao) Update(ctx context.Context, record Row, opt ...Option) (bool, error) {
	res, err := m.model.Update(db.Record(record), m.transOptions(opt...)...)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (m *mysqldao) Delete(ctx context.Context, opt ...Option) (bool, error) {
	res, err := m.model.Delete(m.transOptions(opt...)...)
	if err != nil {
		return false, err
	}
	return res, nil
}
