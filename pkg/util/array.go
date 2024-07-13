package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type ArrStr []string

func (a ArrStr) Has(str string) bool {
	for _, v := range a {
		if v == str {
			return true
		}
	}
	return false
}

func (a ArrStr) Filter(f func(index int, str string) bool) ArrStr {
	var _strs []string
	for i, v := range a {
		if f(i, v) {
			_strs = append(_strs, v)
		}
	}
	return _strs
}

func (a ArrStr) RemoveLast() ArrStr {
	return a.Filter(func(index int, str string) bool {
		return index+1 != len(a)
	})
}

func (a ArrStr) Map(f func(str string) string) ArrStr {
	var _strs []string
	for _, v := range a {
		_strs = append(_strs, f(v))
	}
	return _strs
}

func (a ArrStr) Join(s string) String {
	return String(strings.Join(a, s))
}

func (a ArrStr) Get(index int) String {
	return String(a[index])
}

func (a ArrStr) Concat(arr ...[]string) ArrStr {
	_arr := a.Raw()
	for _, v := range arr {
		_arr = append(_arr, v...)
	}
	return _arr
}

func (a ArrStr) Last() String {
	return String(a.Raw()[len(a)-1])
}

func (a ArrStr) First() String {
	return String(a.Raw()[0])
}

func (a ArrStr) Length() int {
	return len(a)
}

func (a ArrStr) ToSliceInterface() []interface{} {
	var tmp []interface{}
	for _, v := range a {
		tmp = append(tmp, v)
	}
	return tmp
}

func (a ArrStr) Unique() ArrStr {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range a {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (a ArrStr) Raw() []string {
	return a
}

type ArrMap []map[string]interface{}

func (a ArrMap) GetColumn(field string) []interface{} {
	var tmp []interface{}
	for _, v := range a {
		if val, ok := v[field]; ok {
			tmp = append(tmp, val)
		}
	}
	if tmp == nil {
		tmp = []interface{}{}
	}
	return tmp
}

func (a ArrMap) Map(fn func(item map[string]interface{}, index int) map[string]interface{}) ArrMap {
	var newMap []map[string]interface{}
	for i, v := range a {
		newMap = append(newMap, fn(v, i))
	}
	return newMap
}

func (a ArrMap) Filter(fn func(item map[string]interface{}, index int) bool) ArrMap {
	var newMap []map[string]interface{}
	for i, v := range a {
		if fn(v, i) {
			newMap = append(newMap, v)
		}
	}
	return newMap
}

func (a ArrMap) Concat(arr []map[string]interface{}) ArrMap {
	var tmp []map[string]interface{} = a
	for _, v := range arr {
		tmp = append(a, v)
	}
	return tmp
}

func (a ArrMap) Length() int {
	return len(a)
}

func (a ArrMap) ToString(order ...string) string {
	if len(order) == 0 || (len(order) == 1 && order[0] == "*") {
		s, _ := json.Marshal(a)
		return string(s)
	}
	var tmp []string
	for _, v := range a {
		tmp = append(tmp, MapStrInterface(v).ToString(order...))
	}
	return "[" + strings.Join(tmp, ",") + "]"
}

type MapStr map[string]string

func (m MapStr) trans() MapStrInterface {
	var tmp MapStrInterface = map[string]interface{}{}
	for k, v := range m {
		tmp[k] = v
	}

	return tmp
}

func (m MapStr) Merge(data map[string]string) MapStr {
	for k, v := range data {
		m[k] = v
	}
	return m
}

func (m MapStr) ToString(order ...string) string {
	return m.trans().ToString(order...)
}

func (m MapStr) Binding(to interface{}) error {
	return m.trans().Binding(to)
}

type MapStrInterface map[string]interface{}

func (m MapStrInterface) Merge(data map[string]interface{}) MapStrInterface {
	for k, v := range data {
		m[k] = v
	}
	return m
}

func (m MapStrInterface) ToString(order ...string) string {
	if len(order) == 0 {
		s, _ := json.Marshal(m)
		return string(s)
	}
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{', '\n'})
	l := len(order)
	for i, k := range order {
		_, err := fmt.Fprintf(buf, "\t\"%s\": \"%v\"", k, m[k])
		if err != nil {
			continue
		}
		if i < l-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.Write([]byte{'}', '\n'})
	return buf.String()
}

func (m MapStrInterface) Binding(to interface{}) error {
	return Binding(m.ToString(), to)
}

type OrderedArrMap struct {
	Data  *ArrMap
	Order []string
}

func (o OrderedArrMap) MarshalJSON() ([]byte, error) {
	s := o.Data.ToString(o.Order...)
	return []byte(s), nil
}

type ArrInt64 []int64

func (a ArrInt64) SliceInterface() []interface{} {
	var tmp []interface{}
	for _, v := range a {
		tmp = append(tmp, v)
	}
	return tmp
}

func (a ArrInt64) Has(val int64) bool {
	for _, v := range a {
		if v == val {
			return true
		}
	}
	return false
}

func (a ArrInt64) Unique() ArrInt64 {
	keys := make(map[int64]bool)
	var list []int64
	for _, entry := range a {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type ArrInt32 []int32

func (a ArrInt32) SliceInterface() []interface{} {
	var tmp []interface{}
	for _, v := range a {
		tmp = append(tmp, v)
	}
	return tmp
}

func ToInterfaceSlice(els interface{}) ([]interface{}, error) {
	if els == nil {
		return []interface{}{}, nil
	}
	v, err := cast.ToIntSliceE(els)
	if err == nil {
		var tmp []interface{}
		for _, i := range v {
			tmp = append(tmp, i)
		}
		return tmp, nil
	}
	v1, err := cast.ToStringSliceE(els)
	if err == nil {
		var tmp []interface{}
		for _, i := range v1 {
			tmp = append(tmp, i)
		}
		return tmp, nil
	}

	return nil, errors.Wrap(err, "ToInterfaceSlice")
}
