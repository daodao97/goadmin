package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/daodao97/goadmin/pkg/db"
)

func findApplicationTomlPath() (string, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 逐级往上级目录寻找，直到找到根目录
	for currentDir != "/" {
		// 检查当前目录是否是go项目根目录（存在go.mod文件）
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// 在当前目录下寻找cmd目录
			cmdDir := filepath.Join(currentDir, "cmd")
			if _, err := os.Stat(cmdDir); err == nil {
				// 在cmd目录下查找application.toml文件
				applicationTomlPath := filepath.Join(cmdDir, "application.toml")
				if _, err := os.Stat(applicationTomlPath); err == nil {
					return applicationTomlPath, nil
				}
			}
		}

		// 向上级目录继续寻找
		currentDir = filepath.Dir(currentDir)
	}

	return "", fmt.Errorf("未找到 application.toml 文件, 请确保在 goadmin 项目内执行该命令")
}

func yesOrNo(tip string, _default bool) (bool, error) {
	_defaultStr := "n"
	if _default {
		_defaultStr = "y"
	}
	prompt := promptui.Prompt{
		Label:   fmt.Sprintf("%s (y/n)", tip),
		Default: _defaultStr,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "y", nil
}

var re = regexp.MustCompile(`([_\-])([a-zA-Z]+)`)

func ToCamelCase(str string) string {
	camel := re.ReplaceAllString(str, " $2")
	c := cases.Title(language.Dutch)
	camel = c.String(camel)
	return strings.Replace(camel, " ", "", -1)
}

func GetDatabaseNameFromDSN(dsn string) (string, error) {
	// 定义一个正则表达式来匹配数据库名称
	regex := regexp.MustCompile(`/([^/?]+)`)

	// 在DSN中查找匹配的字符串
	matches := regex.FindStringSubmatch(dsn)

	if len(matches) != 2 {
		return "", fmt.Errorf("无法从DSN中提取数据库名称")
	}

	// 返回匹配到的数据库名称
	return matches[1], nil
}

func menuName() (string, error) {
	prompt := promptui.Prompt{
		Label: "菜单名称",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("菜单名称是必须的")
			}
			return nil
		},
	}

	result, err := prompt.Run()

	return result, err
}

func menuDir(m db.Model) (int, error) {
	list := m.Select(db.WhereEq("module_id", 1), db.WhereEq("status", 1), db.WhereEq("type", 1))
	dirs := map[string]int{}
	if list.Err != nil {
		return 0, list.Err
	}
	var items []string
	for _, v := range list.List {
		dirs[cast.ToString(v.Data["name"])] = cast.ToInt(v.Data["id"])
		items = append(items, cast.ToString(v.Data["name"]))
	}
	var id int

	prompt := promptui.SelectWithAdd{
		Label:    "请选择菜单所在目录",
		Items:    items,
		AddLabel: "添加目录",
	}

	index, result, err := prompt.Run()

	if index == -1 {
		_id, err := m.Insert(db.Record{
			"pid":       0,
			"name":      result,
			"type":      1,
			"path":      "#",
			"status":    1,
			"module_id": 1,
		})
		if err != nil {
			return 0, err
		}
		id = cast.ToInt(_id)
	} else {
		for k, v := range dirs {
			if result == k {
				id = v
			}
		}
	}

	return id, err
}
