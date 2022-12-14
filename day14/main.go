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
	max := pathes[0][0].From
	for _, p := range pathes {
		for _, c := range p {
			for _, coord := range []Coord{c.From, c.To} {
				if coord.Y > max.Y {
					max.Y = coord.Y
				}
				if coord.X > max.X {
					max.X = coord.X
				}
			}
		}
	}
	return Coord{Y: 0, X: 0}, Coord{Y: max.Y + 1, X: max.X + 500}
}

type Substance byte

const (
	stone = Substance('#')
	air   = Substance('.')
	sand  = Substance('o')
)

func (s Substance) Passable() bool {
	switch s {
	case air:
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

func makeSpace(min, max Coord, withFloor bool) Space {
	m := make([][]Substance, max.Y-min.Y)
	for i := range m {
		m[i] = make([]Substance, max.X-min.X)
		common.FillSlice(m[i], '.')
	}
	if withFloor {
		line1 := make([]Substance, max.X-min.X)
		line2 := make([]Substance, max.X-min.X)
		common.FillSlice(line1, air)
		common.FillSlice(line2, stone)
		m = append(m, line1, line2)
		max.Y += 2
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

func (s Space) FallSand() bool {
	pos := sandStartPoint
	for {
		if pos.Y+1 >= s.max.Y {
			return false
		}
		if !s.get(pos.Y+1, pos.X).Passable() {
			if s.get(pos.Y+1, pos.X-1).Passable() {
				pos.X--
			} else if s.get(pos.Y+1, pos.X+1).Passable() {
				pos.X++
			} else {
				s.set(pos.Y, pos.X, sand)
				return pos != sandStartPoint
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
		i            int
	)
	for sc.Scan() {
		stonesPathes = append(stonesPathes, parsePath(sc.Text()))
	}
	min, max := findMinMax(stonesPathes)
	// part 1
	space := makeSpace(min, max, false)
	space.FillStones(stonesPathes)

	for i = 0; space.FallSand(); i++ {
	}
	fmt.Println(i)

	// part 2
	space = makeSpace(min, max, true)
	space.FillStones(stonesPathes)
	for i = 0; space.FallSand(); i++ {
	}
	fmt.Println(i + 1)

	// 1016
	// 25402
}
