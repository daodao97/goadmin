package scaffold

import (
	"github.com/go-playground/validator/v10"

	"github.com/daodao97/goadmin/pkg/util"

	"github.com/daodao97/goadmin/pkg/db"
)

const uniqueTokenLessNum = 2
const uniqueTokenFullNum = 3

var unique = util.CustomValidateFunc{
	Handle: func(fl validator.FieldLevel) bool {
		p := fl.Param()
		v := fl.Field().Interface()
		match := util.String(p).Split(".").Raw()
		if len(match) < uniqueTokenLessNum {
			return false
		}
		conn := "default"
		table := ""
		field := ""
		if len(match) == uniqueTokenLessNum {
			table = match[0]
			field = match[1]
		}
		if len(match) == uniqueTokenFullNum {
			conn = match[0]
			table = match[1]
			field = match[2]
		}

		c, err := db.New(table, db.WithConn(conn)).Count(db.WhereEq(field, v))

		if err != nil {
			return false
		}

		return c == 0
	},
	TagName: "unique",
	Message: "数据重复",
}
