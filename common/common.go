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

func ParseSlice[T any](s string, sep string, f func(string) (T, error)) []T {
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

func ParseInt[T int | uint64 | int64 | int32 | uint32 | uint8 | uint16](s string) T {
	num, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		panic("can not parse " + s + ": " + err.Error())
	}
	return T(num)
}

func Keys[K comparable, V any](m map[K]V) []K {
	var res []K
	for k := range m {
		res = append(res, k)
	}
	return res
}

func FillSlice[T any](s []T, el T) {
	for i := 0; i < len(s); i++ {
		s[i] = el
	}
}
