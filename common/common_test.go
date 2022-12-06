package common

import (
	"sort"
	"testing"
)

func TestSortable(t *testing.T) {
	s := []rune("beagfcd")
	sort.Sort(Sortable[rune](s))
	if string(s) != "abcdefg" {
		t.FailNow()
	}
}
