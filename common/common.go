package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
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
	rune | byte | int | float64 | float32 | int64
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

type ReverseSortable[T Comparable] []T

func (s ReverseSortable[T]) Less(i, j int) bool {
	return s[i] > s[j]
}

func (s ReverseSortable[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ReverseSortable[T]) Len() int {
	return len(s)
}

func Abs[T Comparable](a T) T {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}

func ParseArray[T any](s string, sep string, f func(string) (T, error)) []T {
	items := strings.Split(s, sep)
	result := make([]T, len(items))
	var err error
	for i := range items {
		result[i], err = f(strings.TrimSpace(items[i]))
		if err != nil {
			panic(err.Error())
		}
	}
	return result
}

func ParseInt(s string) uint64 {
	num, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		panic("can not parse " + s + ": " + err.Error())
	}
	return uint64(num)
}
