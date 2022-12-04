package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type section struct {
	start, end int
}

func parseLine(line string) [2]section {
	var (
		result [2]section
		err    error
	)
	pairs := strings.Split(line, ",")
	if len(pairs) != 2 {
		panic(fmt.Sprintf("bad line: %s", line))
	}
	for i, p := range pairs {
		sections := strings.Split(p, "-")
		result[i].start, err = strconv.Atoi(sections[0])
		if err != nil {
			panic(fmt.Sprintf("bad line: %s", line))
		}
		result[i].end, err = strconv.Atoi(sections[1])
		if err != nil {
			panic(fmt.Sprintf("bad line: %s", line))
		}
	}
	return result
}

func isSubSecion(s1, s2 section) bool {
	return s1.start <= s2.start && s1.end >= s2.end
}

func isOverlap(s1, s2 section) bool {
	return s1.start <= s2.start && s1.end >= s2.start
}

func main() {
	sc, closeFile := common.FileScanner("./day4/input.txt")
	defer closeFile()

	var cnt1, cnt2 int
	for sc.Scan() {
		line := sc.Text()
		pair := parseLine(line)
		if isSubSecion(pair[0], pair[1]) || isSubSecion(pair[1], pair[0]) {
			cnt1++
		}
		if isOverlap(pair[0], pair[1]) || isOverlap(pair[1], pair[0]) {
			cnt2++
		}
	}
	fmt.Println(cnt1, cnt2)
}
