package util

import (
	"github.com/spf13/cast"
)

const number1 = 1
const number2 = 2
const number3 = 3

type TreeData = map[string]interface{}

// Tree 列表转树形
func Tree(list []TreeData, rootPid int, fields ...string) []TreeData {
	pidKey := "pid"
	idKey := "id"
	childrenKey := "children"
	if len(fields) == number1 {
		pidKey = fields[number1-1]
	}
	if len(fields) == number2 {
		pidKey = fields[number1-1]
		idKey = fields[number2-1]
	}
	if len(fields) == number3 {
		pidKey = fields[number1-1]
		idKey = fields[number2-1]
		childrenKey = fields[number3-1]
	}
	var t []TreeData
	for _, v := range list {
		pid := cast.ToInt(v[pidKey])
		id := cast.ToInt(v[idKey])
		if pid == rootPid {
			child := Tree(list, id, fields...)
			node := v
			if child != nil {
				node[childrenKey] = child
			}
			t = append(t, node)
		}
	}
	return t
}

type TreeConf struct {
	RootPid     int64
	IdField     string
	PidField    string
	ChildrenKey string
	NodesHook   func([]TreeData) []TreeData
}

type TreeOption = func(c *TreeConf)

func treeConf(options ...TreeOption) *TreeConf {
	c := &TreeConf{
		RootPid:     0,
		IdField:     "id",
		PidField:    "pid",
		ChildrenKey: "children",
		NodesHook: func(data []TreeData) []TreeData {
			return data
		},
	}
	for _, o := range options {
		o(c)
	}
	return c
}

func WithPid(pid int64) TreeOption {
	return func(c *TreeConf) {
		c.RootPid = pid
	}
}

func WithIdField(field string) TreeOption {
	return func(c *TreeConf) {
		c.IdField = field
	}
}

func WithPidField(field string) TreeOption {
	return func(c *TreeConf) {
		c.PidField = field
	}
}

func WithChildren(field string) TreeOption {
	return func(c *TreeConf) {
		c.ChildrenKey = field
	}
}

func WithHook(fn func([]TreeData) []TreeData) TreeOption {
	return func(c *TreeConf) {
		c.NodesHook = fn
	}
}

func Tree2(list []TreeData, options ...TreeOption) (t []TreeData) {
	c := treeConf(options...)
	for _, v := range list {
		pid := cast.ToInt64(v[c.PidField])
		id := cast.ToInt64(v[c.IdField])
		if pid == c.RootPid {
			_options := options
			_options = append(_options, WithPid(id))
			child := Tree2(list, _options...)
			node := v
			if child != nil {
				node[c.ChildrenKey] = c.NodesHook(child)
			}
			t = append(t, node)
		}
	}
	t = c.NodesHook(t)
	return t
}
