package common

import (
	"bufio"
	"os"
)

func FileScanner(path string) (*bufio.Scanner, func()) {
	f, err := os.Open(path)
	if err != nil {
		panic("can not open file")
	}
	sc := bufio.NewScanner(f)
	return sc, func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}
}

type Comparable interface {
	rune | byte | int
}

type Sortable[T Comparable] []T

func (s Sortable[T]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Sortable[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sortable[T]) Len() int {
	return len(s)
}
