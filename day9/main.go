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

type CoordSet map[Coord]bool

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

type Knots struct {
	headPos, tailPos Coord
	visitedCoords    CoordSet
}

func NewKnots() Knots {
	return Knots{
		visitedCoords: make(map[Coord]bool),
	}
}

func (k *Knots) MoveHead(move Move) {
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
			k.headPos.X += sign
		} else {
			k.headPos.Y += sign
		}
		k.followHead()
		k.visitedCoords[k.tailPos] = true
		fmt.Println(k.headPos.Y, k.headPos.X, " | ", k.tailPos.Y, k.tailPos.X)
	}
}

func (k *Knots) followHead() {
	if k.headPos.Y > k.tailPos.Y {
		if abs(k.headPos.Y-k.tailPos.Y) > 1 {
			k.tailPos.Y += 1
			k.tailPos.X = k.headPos.X
		}
	} else if k.headPos.Y < k.tailPos.Y {
		if abs(k.tailPos.Y-k.headPos.Y) > 1 {
			k.tailPos.Y -= 1
			k.tailPos.X = k.headPos.X
		}
	}
	if k.headPos.X > k.tailPos.X {
		if abs(k.headPos.X-k.tailPos.X) > 1 {
			k.tailPos.X += 1
			k.tailPos.Y = k.headPos.Y
		}
	} else if k.headPos.X < k.tailPos.X {
		if abs(k.tailPos.X-k.headPos.X) > 1 {
			k.tailPos.X -= 1
			k.tailPos.Y = k.headPos.Y
		}
	}
}

func main() {
	sc, closeFile := common.FileScanner("./day9/input.txt")
	defer closeFile()
	knots := NewKnots()
	for sc.Scan() {
		move := parseMove(sc.Text())
		fmt.Println(sc.Text())
		knots.MoveHead(move)
	}
	fmt.Println(len(knots.visitedCoords))
}
