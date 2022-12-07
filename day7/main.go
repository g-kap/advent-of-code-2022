package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

var fileSizeRe = regexp.MustCompile(`\d+`)

func getFileSize(path string) (int, bool) {
	sizeStr := fileSizeRe.FindString(path)
	if len(sizeStr) > 0 {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			panic(err.Error())
		}
		return size, true
	}
	return 0, false
}

const (
	cdCmdPrefix     = "$ cd "
	lsCmdPrefix     = "$ ls"
	dirOutputPrefix = "dir"
)

func processLine(line string) {
	if strings.Contains(line, cdCmdPrefix) {
		processCDCmd(strings.Replace(line, cdCmdPrefix, "", 1))
	} else if strings.Contains(line, dirOutputPrefix) {

	} else if strings.Contains(line, lsCmdPrefix) {

	} else if fSize, ok := getFileSize(line); ok {
		processFileSize(fSize)
	} else {
		panic(fmt.Sprintf("bad command: %s", line))
	}
}

var (
	curPath []string
	dirSize = make(map[string]int)
)

func processCDCmd(arg string) {
	if arg == ".." {
		curPath = curPath[:len(curPath)-1]
	} else {
		curPath = append(curPath, arg)
	}
}

func processFileSize(fSize int) {
	for i := range curPath {
		dirPath := strings.Join(curPath[:i+1], "/")
		dirSize[dirPath] += fSize
	}
}

func resolvePart1() int {
	var total int
	for _, v := range dirSize {
		if v <= 100000 {
			total += v
		}
	}
	return total
}

func resolvePart2() int {
	totalSpace := 70000000
	freeSpace := totalSpace - dirSize["/"]
	needSpace := 30000000 - freeSpace
	var sizes []int
	for _, v := range dirSize {
		sizes = append(sizes, v)
	}
	sort.Sort(sort.IntSlice(sizes))
	for _, s := range sizes {
		if s >= needSpace {
			return s
		}
	}
	panic("not found")
}

func main() {
	sc, closeFile := common.FileScanner("./day7/input.txt")
	defer closeFile()
	for sc.Scan() {
		processLine(sc.Text())
	}
	fmt.Println(resolvePart1())
	fmt.Println(resolvePart2())
}
