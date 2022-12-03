package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/g-kap/advent-of-code-2022/common"
)

type elfCalories int64

const topElfsCount = 3

type TopElfs struct {
	elfs       [topElfsCount]elfCalories
	weakestIdx int
	sum        elfCalories
}

func (e *TopElfs) compareAndSwapWithWeakest(elf elfCalories) {
	if elf <= e.elfs[e.weakestIdx] {
		return
	}
	e.sum = e.sum - e.elfs[e.weakestIdx] + elf
	e.elfs[e.weakestIdx] = elf
	var newWeakestIdx int
	for i := range e.elfs {
		if e.elfs[i] < e.elfs[newWeakestIdx] {
			newWeakestIdx = i
		}
	}
	e.weakestIdx = newWeakestIdx
}

func main() {
	sc, closeFile := common.FileScanner("./day1/input.txt")
	defer closeFile()
	var (
		lastElf elfCalories
		topElfs TopElfs
	)
	for sc.Scan() {
		caloriesStr := sc.Text()
		if caloriesStr == "" {
			topElfs.compareAndSwapWithWeakest(lastElf)
			lastElf = 0
			continue
		}
		caloriesInt, err := strconv.Atoi(caloriesStr)
		if err != nil {
			log.Fatalf("bad line: %s", caloriesStr)
		}
		lastElf += elfCalories(caloriesInt)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	topElfs.compareAndSwapWithWeakest(lastElf)
	fmt.Println(topElfs.sum)
}
