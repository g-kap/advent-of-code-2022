package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type stack []byte

const bytesPerBox = 4

func parseStackLine(stacks []stack, line []byte) {
	for i, r := range line {
		stackIdx := i / bytesPerBox
		if i%bytesPerBox == 1 {
			stacks[stackIdx] = append(stacks[stackIdx], r)
		}
	}
}

func reverseStacks(stacks []stack) {
	for _, s := range stacks {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
	}
}

func trimStacks(stacks []stack) {
	for i, s := range stacks {
		stacks[i] = []byte(strings.TrimRight(string(s), " "))
	}
}

var re = regexp.MustCompile(`\d+`)

type turn struct {
	count, from, to int
}

func parseTurn(line string) turn {
	nums := re.FindAllString(line, 3)
	if len(nums) != 3 {
		panic("bad line")
	}
	var t turn
	for i, n := range nums {
		intNum, err := strconv.Atoi(n)
		if err != nil {
			panic("bad line")
		}
		switch i {
		case 0:
			t.count = intNum
		case 1:
			t.from = intNum - 1
		case 2:
			t.to = intNum - 1
		}
	}
	return t
}

func makeTurn(stacks []stack, t turn, part2 bool) {
	count := t.count
	if len(stacks[t.from]) < count {
		count = len(stacks[t.from])
	}
	if count == 0 {
		return
	}
	if part2 {
		stacks[t.to] = append(stacks[t.to], stacks[t.from][len(stacks[t.from])-count:]...)
	} else {
		for i := 0; i < count; i++ {
			stacks[t.to] = append(stacks[t.to], stacks[t.from][len(stacks[t.from])-1-i])
		}
	}
	stacks[t.from] = stacks[t.from][:len(stacks[t.from])-count]
}

func main() {
	sc, closeFile := common.FileScanner("./day5/input.txt")
	defer closeFile()

	var (
		lineIdx int = -1
		stacks  []stack
	)
	for sc.Scan() {
		lineIdx++
		line := sc.Bytes()
		var size int
		if lineIdx == 0 {
			size = (len(line) + 1) / bytesPerBox
			stacks = make([]stack, size)
		}
		if string(line) == "" {
			break
		}
		if line[1] == '1' {
			continue
		}
		parseStackLine(stacks, line)
	}
	reverseStacks(stacks)
	trimStacks(stacks)
	for sc.Scan() {
		lineIdx++
		t := parseTurn(sc.Text())
		makeTurn(stacks, t, true)
		trimStacks(stacks)
	}
	for _, s := range stacks {
		fmt.Print(string(s[len(s)-1]))
	}
}
