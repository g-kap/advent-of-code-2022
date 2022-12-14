package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

func init() {
	log.SetFlags(0)
}

func DayFileScanner(dayNum int, example bool) (*bufio.Scanner, func()) {
	suffix := ""
	if example {
		suffix = ".example"
	} else {
		log.SetOutput(io.Discard)
	}
	s := fmt.Sprintf("./day%2d/input%s.txt", dayNum, suffix)
	return FileScanner(s)
}

type ComparableBase interface {
	rune | byte | int | float64 | float32 | int64
}

type Addable interface {
	int | int32 | int64 | float32 | float64 | byte
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

func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	res := make([]T2, len(s))
	for i := range s {
		res[i] = f(s[i])
	}
	return res
}

func Sum[T Addable](s []T) T {
	sum := T(0)
	for i := range s {
		sum += s[i]
	}
	return sum
}

func Contains[T comparable](s []T, v T) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}
	return false
}

func BackMap[T1, T2 comparable](m map[T1]T2) map[T2][]T1 {
	m2 := map[T2][]T1{}
	for k, v := range m {
		m2[v] = append(m2[v], k)
	}
	return m2
}

func Permutations[T comparable](arr []T) [][]T {
	var perm func([]T, int)
	res := [][]T{}

	perm = func(arr []T, n int) {
		if n == 1 {
			tmp := make([]T, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				perm(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	perm(arr, len(arr))
	return res
}

func PermutationsToChan[T comparable](arr []T) <-chan []T {

	ch := make(chan []T, 1)
	var perm func([]T, int)
	perm = func(arr []T, n int) {
		if n == 1 {
			tmp := make([]T, len(arr))
			copy(tmp, arr)
			ch <- tmp
		} else {
			for i := 0; i < n; i++ {
				perm(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	go func() {
		perm(arr, len(arr))
		close(ch)
	}()
	return ch
}

type Set[T comparable] map[T]bool

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Add(v T) {
	s[v] = true
}

func MakeSet[T comparable](s []T) Set[T] {
	res := map[T]bool{}
	for _, v := range s {
		res[v] = true
	}
	return res
}

func Max[T ComparableBase](s ...T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	max := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] > max {
			max = s[i]
		}
	}
	return max
}

func Min[T ComparableBase](s ...T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	min := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] < min {
			min = s[i]
		}
	}
	return min
}
