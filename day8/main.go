package main

import (
	"fmt"

	"github.com/g-kap/advent-of-code-2022/common"
)

func toNumbers(line []byte) []uint8 {
	res := make([]uint8, len(line))
	for i := range res {
		res[i] = line[i] - '0'
	}
	return res
}

func processTrees(m TreeMatrix) {
	var cnt, maxScore int
	for i := range m {
		for j := range m[i] {
			tree := m[i][j]
			left, right, up, down :=
				m.left(i, j), m.right(i, j), m.up(i, j), m.down(i, j)
			if isTallest(tree, left) || isTallest(tree, right) ||
				isTallest(tree, up) || isTallest(tree, down) {
				cnt++
			}
			score := distance(tree, up, true) *
				distance(tree, down, false) *
				distance(tree, right, false) *
				distance(tree, left, true)
			if score > maxScore {
				maxScore = score
			}

		}
	}
	fmt.Println(cnt)
	fmt.Println(maxScore)
}

func isTallest(tree byte, line []byte) bool {
	for _, t := range line {
		if t >= tree {
			return false
		}
	}
	return true
}

func distance(tree byte, line []byte, opposite bool) int {
	var dist int
	if opposite {
		for i := len(line) - 1; i >= 0; i-- {
			dist++
			if line[i] >= tree {
				break
			}
		}
	} else {
		for i := 0; i < len(line); i++ {
			dist++
			if line[i] >= tree {
				break
			}
		}
	}
	return dist
}

type TreeMatrix [][]byte

func (m TreeMatrix) left(i, j int) []byte {
	return m[i][0:j]
}

func (m TreeMatrix) right(i, j int) []byte {
	return m[i][j+1:]
}

func (m TreeMatrix) up(i, j int) []byte {
	s := make([]byte, i)
	for idx := 0; idx < i; idx++ {
		s[idx] = m[idx][j]
	}
	return s
}
func (m TreeMatrix) down(i, j int) []byte {
	s := make([]byte, len(m)-i-1)
	for idx := i + 1; idx < len(m); idx++ {
		s[idx-i-1] = m[idx][j]
	}
	return s
}

func main() {
	sc, closeFile := common.FileScanner("./day8/input.txt")
	defer closeFile()
	var (
		treesMatrix TreeMatrix
	)
	for sc.Scan() {
		treesMatrix = append(treesMatrix, toNumbers(sc.Bytes()))
	}
	processTrees(treesMatrix)
}
