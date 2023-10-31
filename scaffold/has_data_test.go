package scaffold

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_explodeHasStr(t *testing.T) {
	token := []string{
		"pool.db.table:id->lid,a",
		"pool.db.table:id->lid,a as a1",
		"pool.db.table:id->lid,a,b,c",
		"db.table:id->lid,a as a1",
		"db.table:id->lid,a,b,c",
	}

	for _, v := range token {
		spew.Dump(explodeHasStr(v))
	}
}
