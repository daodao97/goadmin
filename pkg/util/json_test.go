package util

import (
	"fmt"
	"testing"
)

func TestUnmarshalWithMacro(t *testing.T) {
	jsonstr := `
{
	"a": "{{ a|raw }}",
	"b": "{{ b|raw }}"
	"c": "{{ .b.jsonstr }}"
}
`
	macro := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"jsonstr": "jsonstr",
			"num":     1,
		},
	}
	fmt.Println(JsonStrVarReplace(jsonstr, macro))
}
