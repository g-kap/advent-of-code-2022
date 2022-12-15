package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Coord struct {
	X, Y int
}

func (c Coord) DistanceTo(c2 Coord) int {
	return common.Abs(c.X-c2.X) + common.Abs(c.Y-c2.Y)
}

type Record struct {
	SensorCoord, BeaconCoord Coord
}

func (r Record) Distance() int {
	return r.BeaconCoord.DistanceTo(r.SensorCoord)
}

func (r Record) Covers(c Coord) bool {
	return r.SensorCoord.DistanceTo(c) <= r.Distance()
}

func findMinMax(records []Record) (Coord, Coord) {
	min, max := records[0].BeaconCoord, records[0].BeaconCoord
	var x, y int
	for _, r := range records {
		x = r.SensorCoord.X - r.Distance()
		if x < min.X {
			min.X = x
		}
		x = r.SensorCoord.X + r.Distance()
		if x > max.X {
			max.X = x
		}
		y = r.SensorCoord.Y - r.Distance()
		if y < min.Y {
			min.Y = y
		}
		y = r.SensorCoord.Y + r.Distance()
		if y > max.Y {
			max.Y = y
		}
	}
	return min, max
}

const (
	sensorPrefix = "Sensor at "
	beaconPrefix = "closest beacon is at "
	sep          = ": "
)

func parseRecord(line string) Record {
	parts := strings.Split(line, sep)
	if len(parts) != 2 {
		panic("bad line")
	}
	sensorCoordsRaw := strings.Split(parts[0][len(sensorPrefix):], ", ")
	beaconCoordsRaw := strings.Split(parts[1][len(beaconPrefix):], ", ")
	return Record{
		SensorCoord: Coord{
			X: common.ParseInt[int](strings.TrimPrefix(sensorCoordsRaw[0], "x=")),
			Y: common.ParseInt[int](strings.TrimPrefix(sensorCoordsRaw[1], "y=")),
		},
		BeaconCoord: Coord{
			X: common.ParseInt[int](strings.TrimPrefix(beaconCoordsRaw[0], "x=")),
			Y: common.ParseInt[int](strings.TrimPrefix(beaconCoordsRaw[1], "y=")),
		},
	}
}

func resolvePart1(records []Record, rowNum int) int {
	min, max := findMinMax(records)
	beaconsMap := make(map[Coord]bool)
	for _, r := range records {
		beaconsMap[r.BeaconCoord] = true
	}
	var sum int
	for x := min.X; x <= max.X; x++ {
	LOOP:
		for _, r := range records {
			c := Coord{X: x, Y: rowNum}
			if !beaconsMap[c] && r.Covers(c) {
				sum++
				break LOOP
			}
		}
	}
	return sum
}

const part2Limit = 4_000_000

//const part2Limit = 20

func minMaxPart2(_ []Record) (Coord, Coord) {
	return Coord{
			X: 0,
			Y: 0,
		}, Coord{
			X: part2Limit,
			Y: part2Limit,
		}
}

type Range struct {
	From, To int
}

func (r Range) Cmp(other Range) int {
	if r.From > other.From {
		return 1
	} else if r.From < other.From {
		return -1
	} else {
		if r.From > other.From {
			return 1
		} else if r.From < other.From {
			return -1
		} else {
			return 0
		}
	}
}

func resolvePart2(records []Record) Coord {
	min, max := minMaxPart2(records)
	beaconsMap := make(map[Coord]bool)
	for _, r := range records {
		beaconsMap[r.BeaconCoord] = true
	}
	for y := min.Y; y <= max.Y; y++ {
		var coveredXRanges []Range
		for _, record := range records {
			offset := common.Abs(record.SensorCoord.Y - y)
			dist := record.Distance() - offset
			if dist > 0 {
				coveredXRanges = append(coveredXRanges, Range{
					record.SensorCoord.X - dist - 1,
					record.SensorCoord.X + dist,
				})
			}
		}
		sort.Sort(common.Sortable2[Range](coveredXRanges))
		to := 0
		for i, r := range coveredXRanges[:len(coveredXRanges)-1] {
			next := coveredXRanges[i+1]
			if r.To > to {
				to = r.To
			}
			if next.To > 0 && to < next.From && r.To < max.X {
				return Coord{X: next.From, Y: y}
			}
		}
	}
	panic("not found")
}

func main() {
	sc, closeFile := common.FileScanner("./day15/input.txt")
	defer closeFile()
	var (
		records []Record
	)
	for sc.Scan() {
		records = append(records, parseRecord(sc.Text()))
	}
	//
	//fmt.Println(resolvePart1(records, 2000000))

	coord := resolvePart2(records)
	fmt.Println(coord, coord.X*part2Limit+coord.Y)
}
