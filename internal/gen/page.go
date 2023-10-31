package gen

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/daodao97/goadmin/pkg/logger"
	"github.com/daodao97/goadmin/pkg/util"
	"github.com/daodao97/goadmin/scaffold"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"os"
)

const (
	columnSql = "select `COLUMN_NAME`, `DATA_TYPE`, `COLUMN_COMMENT` from information_schema.COLUMNS where `TABLE_SCHEMA` = ? and `TABLE_NAME` = ? order by ORDINAL_POSITION"
)

func init() {
	logger.SetLoggerLevel(logger.LevelError)
}

func CratePage(routeName, conn, database, table string) error {
	confFile, err := findApplicationTomlPath()
	if err != nil {
		return err
	}
	conf, err := getConf(confFile)
	if err != nil {
		return err
	}

	dbmap := conf.GetDBMap()
	err = db.Init(dbmap)
	if err != nil {
		return err
	}

	if database == "" {
		database, err = GetDatabaseNameFromDSN(dbmap["default"].DSN)
		if err != nil {
			return err
		}
	}

	fields, err := columns(conn, database, table)

	if len(fields) == 0 {
		return errors.New(fmt.Sprintf("无法获取到表 %s 的字段信息", table))
	}

	pageM := db.New("page", db.ColumnHook(db.Json("page_schema")))

	menuParent, err := menuDir(pageM)
	if err != nil {
		return err
	}
	_menuName, err := menuName()
	if err != nil {
		return err
	}

	err = createPage(routeName, pageM, menuParent, _menuName, fields)
	if err != nil {
		return err
	}
	return nil
}

func getConf(confFile string) (*scaffold.Conf, error) {
	content, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	var c scaffold.Conf
	_, err = toml.Decode(string(content), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type TableField struct {
	Name    string
	Field   string
	Type    string
	Comment string
}

func columns(conn, database, table string) ([]TableField, error) {
	_db, err := db.DB(conn)
	if err != nil {
		return nil, err
	}
	rows, err := _db.Query(columnSql, database, table)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	var fieldList []TableField

	for rows.Next() {
		dest := []interface{}{
			new(string),
			new(string),
			new(string),
		}
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		row := TableField{
			Name:    ToCamelCase(*dest[0].(*string)),
			Field:   *dest[0].(*string),
			Type:    getType(*dest[1].(*string)),
			Comment: getType(*dest[2].(*string)),
		}

		fieldList = append(fieldList, row)
	}
	return fieldList, nil
}

func getType(str string) string {
	switch str {
	case "varchar", "char", "tinytext", "datetime", "text", "longtext", "timestamp":
		return "string"
	case "tinyint":
		return "int"
	default:
		return str
	}
}

func makeSchema(table string, fields []TableField) *scaffold.Schema {
	schema := &scaffold.Schema{
		FormItems: []scaffold.FormItems{},
		Headers:   []scaffold.Header{},
		Filter:    []scaffold.Filter{},
		NormalButton: []scaffold.Button{
			{
				Type:   "jump",
				Text:   "新建",
				Target: fmt.Sprintf("/%s/form", table),
				Props: map[string]interface{}{
					"type": "success",
				},
			},
		},
		RowButton: []scaffold.Button{
			{
				Type:   "jump",
				Text:   "编辑",
				Target: fmt.Sprintf("/%s/{id}", table),
				Props: map[string]interface{}{
					"type": "primary",
				},
			},
		},
		BatchButton: []scaffold.Button{
			{
				Type:   "api",
				Text:   "批量删除",
				Target: fmt.Sprintf("/%s/del", table),
				Props: map[string]interface{}{
					"type": "danger",
				},
				Extra: map[string]interface{}{
					"method": "DELETE",
				},
			},
		},
	}
	ignoreFields := util.ArrStr([]string{"ctime", "mtime", "is_deleted"})
	for _, v := range fields {
		if ignoreFields.Has(v.Field) {
			continue
		}
		label := v.Comment
		if label == "" {
			label = v.Field
		}
		schema.FormItems = append(schema.FormItems, scaffold.FormItems{
			Type:  "input",
			Label: label,
			Field: v.Field,
		})

		schema.Headers = append(schema.Headers, scaffold.Header{
			Label: label,
			Field: v.Field,
		})
	}
	return schema
}

func createPage(tableName string, m db.Model, pid int, name string, fields []TableField) error {
	path := fmt.Sprintf("/%s", tableName)

	row := m.SelectOne(db.WhereEq("path", path))
	if row.Err != nil && !errors.Is(row.Err, db.ErrNotFound) {
		return row.Err
	}

	schema := makeSchema(tableName, fields)

	pageRow := db.Record{
		"pid":         pid,
		"name":        name,
		"path":        fmt.Sprintf("/%s", tableName),
		"type":        2, // 菜单
		"status":      1,
		"module_id":   1,
		"page_schema": schema,
		"page_type":   7, // 通用CRUD
	}

	if row.Data != nil {
		force, err := yesOrNo("该路由已经存在, 确认要覆盖吗?", false)
		if err != nil {
			return err
		}
		if !force {
			return fmt.Errorf("命令终止")
		}
		_, err = m.Update(pageRow, db.WhereEq("id", row.Data["id"]))

		if err != nil {
			return err
		}
	} else {
		_, err := m.Insert(pageRow)
		if err != nil {
			return err
		}
	}

	return nil
}
