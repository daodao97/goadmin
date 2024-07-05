package api

import (
	"github.com/daodao97/goadmin/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

func parseFilters(queryParams map[string][]string) []db.Option {
	var filters []db.Option
	for key, values := range queryParams {
		if strings.HasPrefix(key, "filter.") {
			parts := strings.SplitN(key[7:], ".", 2)
			field := parts[0]
			operator := "eq"
			if len(parts) > 1 {
				operator = parts[1]
			}
			for _, value := range values {
				filters = append(filters, db.Where(field, translateOperator(operator), value))
			}
		}
	}
	return filters
}

func translateOperator(operator string) string {
	switch operator {
	case "eq":
		return "="
	case "gt":
		return ">"
	case "lt":
		return "<"
	case "ge":
		return ">="
	case "le":
		return "<="
	case "like":
		return "like"
	case "not_like":
		return "not_like"
	default:
		return operator
	}
}

func List(c *gin.Context) {
	table := c.Param("table_name")

	ps := c.DefaultQuery("_ps", "20")
	pn := c.DefaultQuery("_pn", "1")

	m := db.New(table)
	opt := parseFilters(c.Request.URL.Query())

	count, err := m.Count(opt...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": map[string]any{
				"total": count,
				"list":  []any{},
				"_ps":   cast.ToInt(ps),
				"_pn":   cast.ToInt(pn),
			},
		})
		return
	}

	opt = append(opt, []db.Option{
		db.Limit(cast.ToInt(ps)),
		db.Offset((cast.ToInt(pn) - 1) * cast.ToInt(ps)),
	}...)

	rows := m.Select(opt...)

	if rows.Err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  rows.Err.Error(),
		})
		return
	}

	list := make([]db.Row, 0)

	if len(rows.List) > 0 {
		list = rows.List
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]any{
			"total": count,
			"list":  list,
			"_ps":   cast.ToInt(ps),
			"_pn":   cast.ToInt(pn),
		},
	})
}

func Create(c *gin.Context) {
	table := c.Param("table_name")
	var requestBody db.Record

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := db.New(table).Insert(requestBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]any{
			"id": id,
		},
	})
}

func Read(c *gin.Context) {
	table := c.Param("table_name")
	id := c.Param("id")

	row := db.New(table).SelectOne(db.WhereEq("id", id))
	if row.Err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  row.Err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": row.Data,
	})
}

func Update(c *gin.Context) {
	table := c.Param("table_name")
	id := c.Param("id")

	var requestBody db.Record

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := db.New(table).UpdateBy(cast.ToInt64(id), requestBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func Delete(c *gin.Context) {
	table := c.Param("table_name")
	id := c.Param("id")

	_, err := db.New(table).Update(db.Record{"is_deleted": 1}, db.WhereEq("id", id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
