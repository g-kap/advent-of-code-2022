package main

import (
	"fmt"
	"regexp"

	"github.com/g-kap/advent-of-code-2022/common"
)

type RobotKind int8

const (
	OreRobot = RobotKind(iota)
	ClayRobot
	ObsidianRobot
	GeodeRobot
)

func (r RobotKind) RR() ResourceRecord {
	rr := ResourceRecord{}
	switch r {
	case OreRobot:
		rr.Ore = 1
	case ClayRobot:
		rr.Clay = 1
	case ObsidianRobot:
		rr.Obsidian = 1
	case GeodeRobot:
		rr.Geode = 1
	default:
		panic("pad resource")
	}
	return rr
}

type Blueprint map[RobotKind]ResourceRecord

type ResourceRecord struct {
	Ore, Clay, Obsidian, Geode int
}

func (rr ResourceRecord) ContainsMoreThan(rr2 ResourceRecord) bool {
	if rr.Ore >= rr2.Ore &&
		rr.Geode >= rr2.Geode &&
		rr.Clay >= rr2.Clay &&
		rr.Obsidian >= rr2.Obsidian {
		return true
	}
	return false
}

func (rr ResourceRecord) add(rr2 ResourceRecord, sign int) ResourceRecord {
	return ResourceRecord{
		Ore:      rr.Ore + (sign * rr2.Ore),
		Clay:     rr.Clay + (sign * rr2.Clay),
		Obsidian: rr.Obsidian + (sign * rr2.Obsidian),
		Geode:    rr.Geode + (sign * rr2.Geode),
	}
}

func (rr ResourceRecord) Add(rr2 ResourceRecord) ResourceRecord {
	return rr.add(rr2, 1)
}

func (rr ResourceRecord) Sub(rr2 ResourceRecord) ResourceRecord {
	return rr.add(rr2, -1)
}

func (rr ResourceRecord) Robots(kind RobotKind) int {
	switch kind {
	case OreRobot:
		return rr.Ore
	case ClayRobot:
		return rr.Clay
	case GeodeRobot:
		return rr.Geode
	case ObsidianRobot:
		return rr.Obsidian
	default:
		panic("bad robot")
	}
}

type cacheRecord struct {
	bpId              int
	resources, robots ResourceRecord
	minute            int
}

var cache = map[cacheRecord]int{}

func setCache(bpId int, resources, robots ResourceRecord, minute int, result int) {
	cache[cacheRecord{
		bpId:      bpId,
		resources: resources,
		robots:    robots,
		minute:    minute,
	}] = result
}

func getCache(bpId int, resources, robots ResourceRecord, minute int) (int, bool) {
	r, ok := cache[cacheRecord{
		bpId:      bpId,
		resources: resources,
		robots:    robots,
		minute:    minute,
	}]
	return r, ok
}

func getBestGeodeAmount(bpId int, bp Blueprint, resources, robots ResourceRecord, minutesLeft int) int {
	if minutesLeft <= 0 {
		return resources.Geode
	}
	if cachedResult, ok := getCache(bpId, resources, robots, minutesLeft); ok {
		return cachedResult
	}
	var results []int

	for _, r := range []RobotKind{GeodeRobot, ObsidianRobot, ClayRobot, OreRobot} {
		if resources.ContainsMoreThan(bp[r]) {
			results = append(results, getBestGeodeAmount(bpId, bp,
				resources.Sub(bp[r]).Add(robots),
				robots.Add(r.RR()),
				minutesLeft-1,
			))
			if r == GeodeRobot {
				break
			}
		}
	}
	results = append(results, getBestGeodeAmount(bpId, bp,
		resources.Add(robots),
		robots,
		minutesLeft-1,
	))
	maxRes := common.Max(results...)
	setCache(bpId, resources, robots, minutesLeft, maxRes)
	return maxRes
}

func main() {
	sc, closeFile := common.FileScanner("./day19/input.example.txt")
	defer closeFile()
	var (
		numRe = regexp.MustCompile(`\d+`)

		blueprints     []Blueprint
		startResources = ResourceRecord{}
		startRobots    = ResourceRecord{Ore: 1}
	)
	for sc.Scan() {
		nums := common.Map[string, int](numRe.FindAllString(sc.Text(), 7), common.ParseInt[int])
		if len(nums) != 7 {
			panic("bad line")
		}
		blueprints = append(blueprints, Blueprint{
			OreRobot:      ResourceRecord{Ore: nums[1]},
			ClayRobot:     ResourceRecord{Ore: nums[2]},
			ObsidianRobot: ResourceRecord{Ore: nums[3], Clay: nums[4]},
			GeodeRobot:    ResourceRecord{Ore: nums[5], Obsidian: nums[6]},
		})
	}

	sum := 0
	for i, bp := range blueprints {
		maxGeodCount := getBestGeodeAmount(i, bp, startResources, startRobots, 24)
		qualityLevel := (i + 1) * maxGeodCount
		sum += qualityLevel
		cache = map[cacheRecord]int{}
	}
	fmt.Println(sum)

	mul := 0
	for i, bp := range blueprints[:common.Min(3, len(blueprints))] {
		maxGeodCount := getBestGeodeAmount(i, bp, startResources, startRobots, 32)
		mul *= maxGeodCount
		cache = map[cacheRecord]int{}
	}
	fmt.Println(mul)
}
