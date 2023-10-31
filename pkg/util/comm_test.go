package util

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"gotest.tools/assert"
)

func TestBinding(t *testing.T) {
	str := `
{
	"a": "1",
	"b": 1,
	"c": [1,2]
}
`
	o := new(struct {
		A int
		B string
		C []int
	})

	err := Binding(str, o)

	assert.Equal(t, err, nil)

	spew.Dump(o)
}
