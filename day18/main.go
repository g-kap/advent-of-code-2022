package main

import (
	"fmt"
	"strconv"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Cube struct {
	X, Y, Z int
}

func isNeighbors(c1, c2 Cube) bool {
	if c1.X == c2.X && c1.Y == c2.Y && common.Abs(c1.Z-c2.Z) == 1 {
		return true
	} else if c1.Y == c2.Y && c1.X == c2.X && common.Abs(c1.Z-c2.Z) == 1 {
		return true
	} else if c1.Z == c2.Z && c1.Y == c2.Y && common.Abs(c1.X-c2.X) == 1 {
		return true
	} else if c1.Y == c2.Y && c1.Z == c2.Z && common.Abs(c1.X-c2.X) == 1 {
		return true
	} else if c1.Z == c2.Z && c1.X == c2.X && common.Abs(c1.Y-c2.Y) == 1 {
		return true
	} else if c1.X == c2.X && c1.Z == c2.Z && common.Abs(c1.Y-c2.Y) == 1 {
		return true
	}
	return false
}

func makeNeighbors(c Cube) []Cube {
	return []Cube{
		{c.X + 1, c.Y, c.Z},
		{c.X, c.Y + 1, c.Z},
		{c.X, c.Y, c.Z + 1},
		{c.X, c.Y, c.Z - 1},
		{c.X, c.Y - 1, c.Z},
		{c.X - 1, c.Y, c.Z},
	}
}

func minMax(lavaCubes []Cube) (Cube, Cube) {
	min, max := lavaCubes[0], lavaCubes[0]
	for _, c := range lavaCubes {
		if c.X < min.X {
			min.X = c.X
		}
		if c.Y < min.Y {
			min.Y = c.Y
		}
		if c.Z < min.Z {
			min.Z = c.Z
		}
		if c.X > max.X {
			max.X = c.X
		}
		if c.Y > max.Y {
			max.Y = c.Y
		}
		if c.Z > max.Z {
			max.Z = c.Z
		}
	}
	min.Z--
	min.Y--
	min.X--
	max.Z++
	max.X++
	max.Y++
	return min, max
}

func makeOpenAirCubes(lavaCubes []Cube) common.Set[Cube] {
	min, max := minMax(lavaCubes)
	cubeSet := common.MakeSet(lavaCubes)
	airCubeSet := common.Set[Cube]{}
	queue := []Cube{min}
	var cube Cube
	for len(queue) > 0 {
		cube, queue = queue[0], queue[1:]
		if cube.X < min.X || cube.Y < min.Y || cube.Z < min.Z || cube.X > max.X || cube.Y > max.Y || cube.Z > max.Z {
			continue
		}
		for _, neighbor := range makeNeighbors(cube) {
			if cubeSet.Contains(neighbor) || airCubeSet.Contains(neighbor) {
				continue
			}
			airCubeSet.Add(neighbor)
			queue = append(queue, neighbor)
		}
	}
	return airCubeSet
}

func resolvePart1(cubes []Cube) int {
	neighbors := map[Cube][]Cube{}
	for _, c1 := range cubes {
		for _, c2 := range cubes {
			if c1 == c2 {
				continue
			}
			if isNeighbors(c1, c2) {
				neighbors[c1] = append(neighbors[c1], c2)
			}
		}
	}
	cnt := 0
	for _, c := range cubes {
		cnt += 6 - len(neighbors[c])
	}
	return cnt
}

func resolvePart2(lavaCubes []Cube) int {
	openAirCubes := makeOpenAirCubes(lavaCubes)
	totalSides := 0
	for _, c := range lavaCubes {
		totalSides += 6
		for _, neighbor := range makeNeighbors(c) {
			if !openAirCubes.Contains(neighbor) {
				totalSides--
			}
		}
	}
	return totalSides
}

func main() {
	sc, closeFile := common.FileScanner("./day18/input.txt")
	defer closeFile()
	var (
		lavaCubes []Cube
	)
	for sc.Scan() {
		coord := common.ParseSlice[int](sc.Text(), ",", strconv.Atoi)
		if len(coord) != 3 {
			panic("bad line")
		}
		lavaCubes = append(lavaCubes, Cube{coord[0], coord[1], coord[2]})
	}
	fmt.Println(resolvePart1(lavaCubes))
	fmt.Println(resolvePart2(lavaCubes))
}
