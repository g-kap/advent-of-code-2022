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

func (m Map) resetScores() {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			m[i][j].Score = -1
		}
	}
}
func (m Map) setBFSScores(from, to *Pos) {
	from.Score = 0
	path := Path{from}
	visited := map[*Pos]bool{}
	var curPos *Pos
	for len(path) > 0 {
		curPos, path = path[0], path[1:]
		if visited[curPos] {
			continue
		}
		visited[curPos] = true
		if curPos == to {
			return
		}
		for _, neighbor := range curPos.Neighbors {
			if !visited[neighbor] {
				neighbor.Score = curPos.Score + 1
				path = append(path, neighbor)
			}
		}
	}
}

func main() {
	sc, closeFile := common.FileScanner("./day12/input.txt")
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
	m.setBFSScores(startPos, destPos)
	fmt.Println(destPos.Score)

	min := destPos.Score
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j].Att == int('a') {
				startPos = m[i][j]
				m.resetScores()
				m.setBFSScores(startPos, destPos)
				if destPos.Score < min && destPos.Score != -1 {
					min = destPos.Score
				}
			}
		}
	}
	fmt.Println(min)
}

//484
//478
