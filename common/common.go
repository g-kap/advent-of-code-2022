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

type ComparableBase interface {
	rune | byte | int | float64 | float32 | int64
}

type Sortable[T ComparableBase] []T

func (s Sortable[T]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Sortable[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sortable[T]) Len() int {
	return len(s)
}

type Comparable[T any] interface {
	Cmp(other T) int
}

type Sortable2[T Comparable[T]] []T

func (s Sortable2[T]) Less(i, j int) bool {
	return s[i].Cmp(s[j]) < 0
}

func (s Sortable2[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sortable2[T]) Len() int {
	return len(s)
}

type ReverseSortable[T ComparableBase] []T

func (s ReverseSortable[T]) Less(i, j int) bool {
	return s[i] > s[j]
}

func (s ReverseSortable[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ReverseSortable[T]) Len() int {
	return len(s)
}

func Abs[T ComparableBase](a T) T {
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

func Keys[K comparable, V any](m map[K]V) []K {
	var res []K
	for k := range m {
		res = append(res, k)
	}
	return res
}
