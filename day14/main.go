package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Coord struct {
	Y, X int
}

type Vector struct {
	From, To Coord
}

type Path []Vector

func parsePath(line string) Path {
	coordsStrSlice := strings.Split(line, "->")
	var (
		coord, prevCoord *Coord
		path             Path
	)
	for _, coordStr := range coordsStrSlice {
		intCoords := common.ParseSlice[int](coordStr, ",", strconv.Atoi)
		if len(intCoords) != 2 {
			panic("bad line")
		}
		coord = &Coord{X: intCoords[0], Y: intCoords[1]}
		if coord != nil && prevCoord != nil {
			if coord.Y != prevCoord.Y && coord.X != prevCoord.X {
				panic("bad line")
			}
			path = append(path, Vector{
				From: *prevCoord,
				To:   *coord,
			})
		}
		prevCoord = &Coord{Y: coord.Y, X: coord.X}
	}
	return path
}

func findMinMax(pathes []Path) (Coord, Coord) {
	min := pathes[0][0].From
	max := pathes[0][0].From

	for _, p := range pathes {
		for _, c := range p {
			for _, coord := range []Coord{c.From, c.To} {
				if coord.Y < min.Y {
					min.Y = coord.Y
				} else if coord.Y > max.Y {
					max.Y = coord.Y
				}
				if coord.X < min.X {
					min.X = coord.X
				} else if coord.X > max.X {
					max.X = coord.X
				}
			}
		}
	}
	min.Y = 0
	min.X--
	max.X += 50
	max.Y++
	return min, max
}

type Substance byte

const (
	stone = Substance('#')
	air   = Substance('.')
	sand  = Substance('o')
	track = Substance('~')
)

func (s Substance) Passable() bool {
	switch s {
	case air, track:
		return true
	case sand, stone:
		return false
	}
	panic("bad Substance")
}

type Space struct {
	min, max Coord
	space    [][]Substance
}

func makeSpace(min, max Coord) Space {
	m := make([][]Substance, max.Y-min.Y)
	for i := range m {
		m[i] = make([]Substance, max.X-min.X)
		common.FillSlice(m[i], '.')
	}
	return Space{
		min:   min,
		max:   max,
		space: m,
	}
}

func (s Space) get(y, x int) Substance {
	y -= s.min.Y
	x -= s.min.X
	return s.space[y][x]
}
func (s Space) set(y, x int, m Substance) {
	y -= s.min.Y
	x -= s.min.X
	s.space[y][x] = m
}

func (s Space) FillStones(stonesPathes []Path) {
	for _, path := range stonesPathes {
		for _, vector := range path {
			if vector.From.X == vector.To.X {
				if vector.From.Y < vector.To.Y {
					for i := vector.From.Y; i <= vector.To.Y; i++ {
						s.space[i-s.min.Y][vector.From.X-s.min.X] = '#'
					}
				} else {
					for i := vector.From.Y; i >= vector.To.Y; i-- {
						s.space[i-s.min.Y][vector.From.X-s.min.X] = '#'
					}
				}
			} else if vector.From.Y == vector.To.Y {
				if vector.From.X < vector.To.X {
					for i := vector.From.X; i <= vector.To.X; i++ {
						s.space[vector.From.Y-s.min.Y][i-s.min.X] = '#'
					}
				} else {
					for i := vector.From.X; i >= vector.To.X; i-- {
						s.space[vector.From.Y-s.min.Y][i-s.min.X] = '#'
					}
				}
			} else {
				panic("bad vector")
			}
		}
	}
}

var sandStartPoint = Coord{
	Y: 0,
	X: 500,
}

func (s Space) FallSand(withTrack bool) bool {
	pos := sandStartPoint
	for {
		if pos.Y+1 >= s.max.Y {
			return false
		}
		if withTrack {
			s.set(pos.Y, pos.X, track)
		}
		if !s.get(pos.Y+1, pos.X).Passable() {
			if s.get(pos.Y+1, pos.X-1).Passable() {
				pos.X--
			} else if s.get(pos.Y+1, pos.X+1).Passable() {
				pos.X++
			} else {
				s.set(pos.Y, pos.X, sand)
				return true
			}
		}
		pos.Y++
	}
}

func main() {
	sc, closeFile := common.FileScanner("./day14/input.txt")
	defer closeFile()
	var (
		stonesPathes []Path
	)
	for sc.Scan() {
		stonesPathes = append(stonesPathes, parsePath(sc.Text()))
	}
	min, max := findMinMax(stonesPathes)
	space := makeSpace(min, max)
	space.FillStones(stonesPathes)
	var i int
	for i = 0; space.FallSand(false); i++ {
	}
	space.FallSand(true)
	printSpace(space)
	fmt.Println(i)
}

func printSpace(s Space) {
	ws := ""
	for line := 0; line < 3; line++ {
		fmt.Print("    |")
		for i := s.min.X; i < s.max.X; i++ {
			fmt.Print(string(fmt.Sprintf("%3d", i)[line]) + ws)
		}
		fmt.Println()
	}
	for i := 0; i < len(s.space); i++ {
		fmt.Printf("%3d |", i)
		for j := 0; j < len(s.space[i]); j++ {
			fmt.Printf(string(s.space[i][j]) + ws)
		}
		fmt.Println()
	}
}
