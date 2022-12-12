package main

import (
	"fmt"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Map [][]*Pos

type Pos struct {
	X, Y      int
	Att       int
	Neighbors []*Pos
	Score     int
}

func (m Map) initNeighbors() {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if i == 3 && j == 6 {
				fmt.Println()
			}
			var potentialNeighbors []*Pos
			// up
			if i > 0 {
				potentialNeighbors = append(potentialNeighbors, m[i-1][j])
			}
			// down
			if i < len(m)-1 {
				potentialNeighbors = append(potentialNeighbors, m[i+1][j])
			}
			// right
			if j > 0 {
				potentialNeighbors = append(potentialNeighbors, m[i][j-1])
			}
			// left
			if j < len(m[i])-1 {
				potentialNeighbors = append(potentialNeighbors, m[i][j+1])
			}
			for _, pn := range potentialNeighbors {
				if pn.Att <= m[i][j].Att+1 {
					m[i][j].Neighbors = append(m[i][j].Neighbors, pn)
				}
			}
		}
	}
}

type Path []*Pos

func (p *Pos) findAllPathes(m Map, dest *Pos) []Path {
	var allPathes []Path
	walkRecursive(m, &allPathes, nil, p, dest)
	return allPathes
}

func walkRecursive(
	m Map,
	pathes *[]Path, path Path, start, dest *Pos,
) {
	path = append(path, start)
	for _, neighbor := range start.Neighbors {
		var visited bool
		for _, pos := range path {
			if neighbor == pos {
				visited = true
				break
			}
		}
		if visited {
			continue
		}
		if neighbor == dest {
			*pathes = append(*pathes, path)
			return
		}
		walkRecursive(m, pathes, path, neighbor, dest)
	}
}

func main() {
	sc, closeFile := common.FileScanner("./day12/input.example.txt")
	defer closeFile()
	var (
		m                 Map
		startPos, destPos *Pos
		i                 int = -1
	)
	for sc.Scan() {
		i++
		line := sc.Bytes()
		m = append(m, make([]*Pos, len(line)))
		for j, b := range line {
			switch b {
			case 'S':
				m[i][j] = &Pos{j, i, int('a'), nil, -1}
				startPos = m[i][j]
			case 'E':
				m[i][j] = &Pos{j, i, int('z'), nil, -1}
				destPos = m[i][j]
			default:
				if b < 'a' || b > 'z' {
					panic("bad line")
				}
				m[i][j] = &Pos{j, i, int(b), nil, -1}
			}
		}
	}
	m.initNeighbors()
	pathes := startPos.findAllPathes(m, destPos)
	var smallestIdx int
	for i := range pathes {
		if len(pathes[i]) < len(pathes[smallestIdx]) {
			smallestIdx = i
		}
	}

	printPath(m, pathes[smallestIdx])
	fmt.Println(len(pathes[smallestIdx]))
}

func printPath(m Map, path Path) {
	fmt.Println("================================================")
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			s := "  . "
			for x, p := range path {
				if p.X == j && p.Y == i {
					s = fmt.Sprintf("%2d%c ", x, m[i][j].Att)
				}
			}
			fmt.Print(s, " ")
		}
		fmt.Println()
	}
	fmt.Println("================================================")
}
