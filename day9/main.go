package main

import "C"
import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Coord struct {
	X, Y int
}

var (
	abs = common.Abs[int]
)

type Move Coord

type CoordSet map[Coord]int

func parseMove(s string) Move {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		panic("bad line")
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("bad line")
	}
	var move Move
	switch parts[0] {
	case "R":
		move.X = num
	case "L":
		move.X = -num
	case "U":
		move.Y = num
	case "D":
		move.Y = -num
	default:
		panic("bad line")
	}
	return move
}

type Knot struct {
	pos           Coord
	next          *Knot
	visitedCoords CoordSet
}

func NewKnot(childrenCount int) *Knot {
	knot := &Knot{
		visitedCoords: make(CoordSet),
	}
	parent := knot
	if childrenCount > 0 {
		for i := 0; i < childrenCount; i++ {
			parent.next = NewKnot(0)
			parent = parent.next
		}
	}
	return knot
}

func (k *Knot) follow(head Coord) {
	if abs(head.Y-k.pos.Y) < 2 && abs(head.X-k.pos.X) < 2 && abs(k.pos.Y-head.Y) < 2 && abs(k.pos.X-head.X) < 2 {
		return
	}
	if k.pos.X == head.X && (abs(k.pos.Y-head.Y) > 1 || abs(head.Y-k.pos.Y) > 1) {
		if head.Y > k.pos.Y {
			k.pos.Y++
		} else if head.Y < k.pos.Y {
			k.pos.Y--
		}
	} else if k.pos.Y == head.Y {
		if head.X > k.pos.X {
			k.pos.X++
		} else if head.X < k.pos.X {
			k.pos.X--
		}
	} else {
		if head.Y > k.pos.Y && head.X > k.pos.X {
			k.pos.Y++
			k.pos.X++
		} else if head.Y > k.pos.Y && head.X < k.pos.X {
			k.pos.Y++
			k.pos.X--
		} else if head.Y < k.pos.Y && head.X > k.pos.X {
			k.pos.Y--
			k.pos.X++
		} else if head.Y < k.pos.Y && head.X < k.pos.X {
			k.pos.Y--
			k.pos.X--
		}
	}
	k.visitedCoords[k.pos] = 1
	if k.next != nil {
		k.next.follow(k.pos)
	}
}

func (k *Knot) MoveHead(move Move) {
	var (
		steps, sign int
	)
	if move.X != 0 {
		steps = int(math.Abs(float64(move.X)))
		sign = move.X / steps
	} else if move.Y != 0 {
		steps = int(math.Abs(float64(move.Y)))
		sign = move.Y / steps
	} else {
		panic("bad move")
	}
	for i := 0; i < steps; i++ {
		if move.X != 0 {
			k.pos.X += sign
		} else {
			k.pos.Y += sign
		}
		k.visitedCoords[k.pos] = 1
		if k.next != nil {
			k.next.follow(k.pos)
		}
	}
}

func (k *Knot) GetTail() *Knot {
	if k.next == nil {
		return k
	} else {
		return k.next.GetTail()
	}
}

func main() {
	sc, closeFile := common.FileScanner("./day9/input.txt")
	defer closeFile()
	knot := NewKnot(9)
	for sc.Scan() {
		move := parseMove(sc.Text())
		knot.MoveHead(move)
	}
	visited := knot.GetTail().visitedCoords
	fmt.Println(len(visited)) // 2372
}
