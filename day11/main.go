package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type WorryLevel uint64

type Monkey struct {
	startingItems []WorryLevel
	operation     MonkeyOp
	test          Test
}

type MonkeyOp func(old WorryLevel) WorryLevel

type Test struct {
	divisibleBy uint64
	ifTrue      uint64
	ifFalse     uint64
}

const (
	startingItemsPrefix = "  Starting items: "
	operationPrefix     = "  Operation: new = old "
	divisibleByPrefix   = "  Test: divisible by "
	ifTruePrefix        = "    If true: throw to monkey "
	ifFalsePrefix       = "    If false: throw to monkey "
)

func parseOp(line string) MonkeyOp {
	parts := strings.Split(strings.TrimPrefix(line, operationPrefix), " ")
	if len(parts) != 2 {
		panic("bad operation: " + line)
	}
	sign, arg := parts[0], parts[1]
	if sign != "+" && sign != "*" {
		panic("bad operation: " + line)
	}
	if arg == "old" && sign == "+" {
		return func(old WorryLevel) WorryLevel {
			return old + old
		}
	} else if arg == "old" && sign == "*" {
		return func(old WorryLevel) WorryLevel {
			return old * old
		}
	}
	argNumInt, err := strconv.ParseUint(arg, 10, 64)
	argNum := WorryLevel(argNumInt)
	if err != nil {
		panic("bad operation: " + line)
	}
	if sign == "+" {
		return func(old WorryLevel) WorryLevel {
			return old + argNum
		}
	} else if sign == "*" {
		return func(old WorryLevel) WorryLevel {
			return old * argNum
		}
	}
	panic("bad operation: " + line)
}

func main() {
	sc, closeFile := common.FileScanner("./day11/input.txt")
	defer closeFile()
	var (
		monkeys   []Monkey
		monkeyIdx = -1
	)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "Monkey") {
			monkeys = append(monkeys, Monkey{})
			monkeyIdx++
		} else if strings.HasPrefix(line, startingItemsPrefix) {
			monkeys[monkeyIdx].startingItems = common.ParseArray[WorryLevel](strings.TrimPrefix(line, startingItemsPrefix), ",",
				func(s string) (WorryLevel, error) { a, _ := strconv.ParseInt(s, 10, 64); return WorryLevel(a), nil })
		} else if strings.HasPrefix(line, operationPrefix) {
			monkeys[monkeyIdx].operation = parseOp(line)
		} else if strings.HasPrefix(line, divisibleByPrefix) {
			monkeys[monkeyIdx].test.divisibleBy = common.ParseInt(line[len(divisibleByPrefix):])
		} else if strings.HasPrefix(line, ifTruePrefix) {
			monkeys[monkeyIdx].test.ifTrue = common.ParseInt(line[len(ifTruePrefix):])
		} else if strings.HasPrefix(line, ifFalsePrefix) {
			monkeys[monkeyIdx].test.ifFalse = common.ParseInt(line[len(ifFalsePrefix):])
		} else if line == "" {
		} else {
			panic("bad line")
		}
	}

	lcm := findLCM(monkeys)

	var worryLevel WorryLevel
	counters := make([]int64, len(monkeys))
	for round := 1; round <= 10_000; round++ {
		for idx, monkey := range monkeys {
			for _, item := range monkey.startingItems {
				counters[idx]++
				worryLevel = WorryLevel(uint64(item) % lcm)
				worryLevel = monkey.operation(worryLevel)
				toMonkey := monkey.test.ifFalse
				if worryLevel%WorryLevel(monkey.test.divisibleBy) == 0 {
					toMonkey = monkey.test.ifTrue
				}
				monkeys[toMonkey].startingItems = append(
					monkeys[toMonkey].startingItems, worryLevel)
			}
			monkeys[idx].startingItems = nil
		}
	}
	sort.Sort(common.ReverseSortable[int64](counters))
	fmt.Println(counters[0] * counters[1])
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func findLCM(monkeys []Monkey) uint64 {
	var intSlice []int
	for _, m := range monkeys {
		intSlice = append(intSlice, int(m.test.divisibleBy))
	}
	return uint64(LCM(intSlice[0], intSlice[1], intSlice[2:]...))
}
