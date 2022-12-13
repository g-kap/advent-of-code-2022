package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Record []any

type Pair struct {
	first, second Record
}

func (r Record) Cmp(other Record) int {
	if len(r) == 0 && len(other) == 0 {
		return 0
	}
	rLen, otherLen := len(r), len(other)
	for i := 0; ; i++ {
		if i >= rLen || i >= otherLen {
			if i >= rLen && i < otherLen {
				return -1
			} else if i >= otherLen && i < rLen {
				return 1
			} else {
				return 0
			}
		}
		rInt, rIsInt := r[i].(float64)
		otherInt, otherIsInt := other[i].(float64)
		if rIsInt && otherIsInt {
			if rInt < otherInt {
				return -1
			} else if rInt > otherInt {
				return 1
			} else {
				continue
			}
		}
		var r2, other2 Record
		if rIsInt {
			r2 = []any{r[i]}
		} else {
			r2 = r[i].([]any)
		}
		if otherIsInt {
			other2 = []any{other[i]}
		} else {
			other2 = other[i].([]any)
		}
		cmp := r2.Cmp(other2)
		if cmp == 0 {
			continue
		} else {
			return cmp
		}
	}
}

var (
	divider2 = Record{[]any{float64(2)}}
	divider6 = Record{[]any{float64(6)}}
)

func main() {
	sc, closeFile := common.FileScanner("./day13/input.txt")
	defer closeFile()
	var (
		pairs         []Pair
		plainList     []Record
		first, second Record
	)
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		if first == nil {
			err := json.Unmarshal(sc.Bytes(), &first)
			if err != nil {
				panic("bad line")
			}
			plainList = append(plainList, first)
		} else if second == nil {
			err := json.Unmarshal(sc.Bytes(), &second)
			if err != nil {
				panic("bad line")
			}
			pairs = append(pairs, Pair{first, second})
			plainList = append(plainList, second)
			first, second = nil, nil
		}
	}
	plainList = append(plainList, divider2, divider6)
	sort.Sort(common.Sortable2[Record](plainList))

	resolvePart1(pairs)
	resolvePart2(plainList)
}

func resolvePart1(pairs []Pair) {
	sum := 0
	for idx, p := range pairs {
		if p.first.Cmp(p.second) == -1 {
			sum += idx + 1
		}
	}
	fmt.Println(sum)
}

func resolvePart2(p []Record) {
	var idx2, idx6 int
	for idx, r := range p {
		if r.Cmp(divider2) == 0 {
			idx2 = idx + 1
		}
		if r.Cmp(divider6) == 0 {
			idx6 = idx + 1
		}
	}
	fmt.Println(idx2 * idx6)
}
