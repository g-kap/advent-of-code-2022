package main

import "C"
import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type WorryLevel *big.Int

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
			return (*big.Int)(old).Add(old, old)
		}
	} else if arg == "old" && sign == "*" {
		return func(old WorryLevel) WorryLevel {
			return (*big.Int)(old).Mul(old, old)
		}
	}
	argNumInt, err := strconv.ParseInt(arg, 10, 64)
	argNum := big.NewInt(argNumInt)
	if err != nil {
		panic("bad operation: " + line)
	}
	if sign == "+" {
		return func(old WorryLevel) WorryLevel {
			//return old + argNum
			return (*big.Int)(old).Add(old, argNum)
		}
	} else if sign == "*" {
		return func(old WorryLevel) WorryLevel {
			//return old * argNum
			return (*big.Int)(old).Mul(old, argNum)
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
				func(s string) (WorryLevel, error) {
					a, _ := strconv.ParseInt(s, 10, 64)
					return big.NewInt(a), nil
				})
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

	var worryLevel WorryLevel
	counters := make([]int64, len(monkeys))
	for round := 1; round <= 10_000; round++ {
		fmt.Println("Round ", round)
		for idx, monkey := range monkeys {
			for _, item := range monkey.startingItems {
				counters[idx]++
				worryLevel = monkey.operation(item)
				//worryLevel = worryLevel / 3
				toMonkey := monkey.test.ifFalse
				//if worryLevel%WorryLevel(monkey.test.divisibleBy) == 0 {
				wl := big.NewInt(0).Set(worryLevel)
				if wl.Mod(worryLevel, big.NewInt(int64(monkey.test.divisibleBy))).Cmp(big.NewInt(0)) == 0 {
					toMonkey = monkey.test.ifTrue
				}
				monkeys[toMonkey].startingItems = append(
					monkeys[toMonkey].startingItems, worryLevel)
			}
			monkeys[idx].startingItems = nil
		}
		prnt(round, counters)
	}
	sort.Sort(common.Sortable[int64](counters))
	fmt.Println(counters)
}

func prnt(r int, args ...any) {
	switch r {
	case 1, 20, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 10000:
		fmt.Println(args...)
	}
}
