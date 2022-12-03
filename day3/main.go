package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func weight(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r - 'a' + 1)
	} else if r >= 'A' && r <= 'Z' {
		return int(r - 'A' + 27)
	} else {
		panic("bad rune")
	}
}

func commonRune(lines []string) rune {
	if len(lines) < 2 {
		panic("bad lines")
	}
	for _, r := range lines[0] {
		ok := false
	LOOP:
		for _, line := range lines[1:] {
			if strings.ContainsRune(line, r) {
				ok = true
			} else {
				ok = false
				break LOOP
			}
		}
		if ok {
			return r
		}
	}
	panic("common rune not found")
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("can not open ./input.txt")
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var (
		totalPart1, totalPart2 int
		lines                  []string
	)
	for sc.Scan() {
		line := sc.Text()
		totalPart1 += weight(commonRune([]string{line[:len(line)/2], line[len(line)/2:]}))
		lines = append(lines, line)
		if len(lines) == 3 {
			totalPart2 += weight(commonRune(lines))
			lines = nil
		}
	}
	fmt.Println(totalPart1)
	fmt.Println(totalPart2)
}
