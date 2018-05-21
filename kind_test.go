package bigjsonvalue

import (
	"testing"
)

func TestKindNames(t *testing.T) {
	for k := Nil; k < lastKind; k++ {
		if k.String() != kindNames[k] {
			t.Errorf("Kind(%d).String() unexpectedly returns <%s>, should be <%s>",
				uint(k), k.String(), kindNames[k])
		}
	}
}
