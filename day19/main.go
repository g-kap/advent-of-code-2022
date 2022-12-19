package main

import (
	"fmt"
	"regexp"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Resources map[ResType]int

func (r Resources) RR() ResourceRecord {
	return ResourceRecord{
		Ore:      r[Ore],
		Clay:     r[Clay],
		Obsidian: r[Obsidian],
		Geode:    r[Geode],
	}
}

type ResType int8

const (
	Ore = ResType(iota)
	Clay
	Obsidian
	Geode
)

type RobotKind ResType

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

const (
	OreRobot      = RobotKind(Ore)
	ClayRobot     = RobotKind(Clay)
	ObsidianRobot = RobotKind(Obsidian)
	GeodeRobot    = RobotKind(Geode)
)

type Blueprint map[RobotKind]Resources

const maxTime = 24

type ResourceRecord struct {
	Ore, Clay, Obsidian, Geode int
}

func (rr ResourceRecord) HasEnoughForBuild(bp Blueprint, r RobotKind) bool {
	if rr.Ore >= bp[r][Ore] &&
		rr.Geode >= bp[r][Geode] &&
		rr.Clay >= bp[r][Clay] &&
		rr.Obsidian >= bp[r][Obsidian] {
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

func getBestGeodeAmount(bp Blueprint, resources, robots ResourceRecord, minute int) int {
	if minute > maxTime {
		return resources.Geode
	}
	var results []int

	for _, r := range []RobotKind{GeodeRobot, ObsidianRobot, ClayRobot, OreRobot} {
		if resources.HasEnoughForBuild(bp, r) {
			results = append(results, getBestGeodeAmount(bp,
				resources.Sub(bp[r].RR()).Add(robots),
				robots.Add(r.RR()),
				minute+1,
			))
			if r == GeodeRobot {
				break
			}
		}
	}
	results = append(results, getBestGeodeAmount(bp,
		resources.Add(robots),
		robots,
		minute+1,
	))

	return common.Max(results...)
}

func main() {
	sc, closeFile := common.FileScanner("./day19/input.txt")
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
			OreRobot:      Resources{Ore: nums[1]},
			ClayRobot:     Resources{Ore: nums[2]},
			ObsidianRobot: Resources{Ore: nums[3], Clay: nums[4]},
			GeodeRobot:    Resources{Ore: nums[5], Obsidian: nums[6]},
		})
	}

	sum := 0
	for i, bp := range blueprints {
		maxGeodCount := getBestGeodeAmount(bp, startResources, startRobots, 1)
		qualityLevel := (i + 1) * maxGeodCount
		fmt.Println(i+1, maxGeodCount, qualityLevel)
		sum += qualityLevel
	}
	fmt.Println(sum)
}
