package cconf

import (
	"github.com/gin-gonic/gin"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/scaffold"
	"github.com/daodao97/goadmin/scaffold/dao"
)

func newService(s *scaffold.Scaffold) *service {
	s.SetModel(db.New("common_config", db.ColumnHook(db.Json("value"))))
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

func (s service) getCommConfById(ctx *gin.Context, id string) (dao.Row, error) {
	d, err := s.Dao(ctx)
	if err != nil {
		return nil, err
	}
	row, err := d.SelectOne(ctx, dao.Where("id", "=", id))
	if err != nil {
		return nil, err
	}
	return row, nil
}
