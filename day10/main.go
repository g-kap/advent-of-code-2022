package main

import "C"
import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/g-kap/advent-of-code-2022/common"
)

type OpType uint8

const (
	opNoop = OpType(iota)
	opAddX
)

type Op struct {
	Type OpType
	Arg  int64
}

func parseOp(s string) Op {
	if strings.HasPrefix(s, "noop") {
		return Op{Type: opNoop}
	} else if strings.HasPrefix(s, "addx") {
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			panic("bad line")
		}
		intArg, err := strconv.Atoi(parts[1])
		if err != nil {
			panic("bad line")
		}
		return Op{Type: opAddX, Arg: int64(intArg)}
	} else {
		panic("bad line")
	}
}

type CPU struct {
	regX     int64
	pipeline chan Op
	counter  int64
	wg       sync.WaitGroup
	crt      [6][40]byte
}

func NewCPU() *CPU {
	c := CPU{
		regX:     1,
		pipeline: make(chan Op, 1),
	}
	c.wg.Add(1)
	go func() {
		var sum int64
		for op := range c.pipeline {
			row := c.counter / 40 % 6
			col := c.counter % 40
			char := byte('.')
			if col >= c.regX-1 && col <= c.regX+1 {
				char = '#'
			}
			c.crt[row][col] = char

			c.counter++
			if c.counter == 20 || (c.counter-20)%40 == 0 {
				pow := c.counter * c.regX
				sum += pow
			}
			switch op.Type {
			case opNoop:
			case opAddX:
				c.regX += op.Arg
			}
		}
		fmt.Println(sum)
		c.wg.Done()
	}()
	return &c
}

func (c *CPU) Apply(op Op) {
	switch op.Type {
	case opNoop:
		c.pipeline <- op
	case opAddX:
		c.pipeline <- Op{Type: opNoop}
		c.pipeline <- op
	default:
		panic("bad op")
	}
}

func (c *CPU) Close() {
	close(c.pipeline)
	c.wg.Wait()
}

func main() {
	sc, closeFile := common.FileScanner("./day10/input.txt")
	defer closeFile()
	var (
		cpu = NewCPU()
	)
	for sc.Scan() {
		op := parseOp(sc.Text())
		cpu.Apply(op)
	}
	cpu.Close()
	for _, row := range cpu.crt {
		fmt.Println(string(row[:]))
	}
}

// 11720
// ERCREPCJ
//####.###...##..###..####.###...##....##.
//#....#..#.#..#.#..#.#....#..#.#..#....#.
//###..#..#.#....#..#.###..#..#.#.......#.
//#....###..#....###..#....###..#.......#.
//#....#.#..#..#.#.#..#....#....#..#.#..#.
//####.#..#..##..#..#.####.#.....##...##..
