package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("can not open ./input.txt")
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
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
