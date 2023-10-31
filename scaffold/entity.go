package scaffold

import (
	"github.com/daodao97/goadmin/pkg/db"
)

func NewUser() db.Model {
	return db.New("user", db.ColumnHook(db.Json("role_ids")))
}

func NewPage() db.Model {
	return db.New("page", db.ColumnHook(db.Json("page_schema"), db.CommaInt("owner_ids")))
}

func NewCommonConfModel() db.Model {
	return db.New("common_config")
}

func NewRole() db.Model {
	return db.New("role")
}

type UserConfig struct {
	RoleType []int64 `json:"role_type"` // 0 无限制, 1 资源只读
}

type role struct {
	Resource []string    `json:"resource"`
	Config   *UserConfig `json:"config"`
}

func (e *role) Table() string {
	return "role"
}
