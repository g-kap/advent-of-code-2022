package main

import (
	"fmt"
	"sort"

	"github.com/g-kap/advent-of-code-2022/common"
)

const (
	startOfPacket  = 4
	startOfMessage = 14
)

func findStartPos(line string, numOfBytes int) int {
	for i := numOfBytes; i < len(line); i++ {
		if !hasDuplicates(line[i-numOfBytes : i]) {
			return i
		}
	}
	panic("not found")
}

func hasDuplicates(line string) bool {
	s := []rune(line)
	sort.Sort(common.Sortable[rune](s))
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			return true
		}
	}
	return false
}

func main() {
	sc, closeFile := common.FileScanner("./day6/input.txt")
	defer closeFile()
	for sc.Scan() {
		fmt.Println(
			findStartPos(sc.Text(), startOfPacket),
			findStartPos(sc.Text(), startOfMessage))
	}
}
