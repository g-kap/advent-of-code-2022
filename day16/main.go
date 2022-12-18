package main

import (
	"fmt"
	"strings"

	"github.com/g-kap/advent-of-code-2022/common"
)

type ValveID [2]byte

func (v ValveID) String() string {
	return string(v[:])
}

func ValveIDFromStr(s string) ValveID {
	var v ValveID
	if len(s) != 2 {
		panic("bad valve id")
	}
	v[0] = s[0]
	v[1] = s[1]
	return v
}

type Valve struct {
	ID       ValveID
	FlowRate int
	LeadsTo  []*Valve
	IsOpen   bool
}

func (v *Valve) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%s: %d -> %s",
		v.ID, v.FlowRate,
		common.Map(v.LeadsTo, func(v *Valve) string { return string(v.ID[:]) }),
	)
}

type ValveMap map[ValveID]*Valve

func (vm ValveMap) TotalOpenedFlowRate() int {
	total := 0
	for _, v := range vm {
		if v.IsOpen {
			total += v.FlowRate
		}
	}
	return total
}

func (vm ValveMap) AllValvesOpened() bool {
	for _, v := range vm {
		if !v.IsOpen && v.FlowRate > 0 {
			return false
		}
	}
	return true
}

func (m ValveMap) Get(id string) *Valve {
	return m[ValveID{id[0], id[1]}]
}

func (m ValveMap) Reset() {
	for _, v := range m {
		v.IsOpen = false
	}
}

func parseValveParams(line string) (ValveID, int, []ValveID) {
	parts := strings.Split(line, "; ")
	if len(parts) != 2 {
		panic("bad line")
	}
	valveIDStr := parts[0][len("Valve ") : len("Valve ")+2]
	flowRateStr := parts[0][len("Valve AA has flow rate="):]
	leadsToPrefLen := len("tunnel lead to valves ")
	if strings.HasPrefix(parts[1], "tunnels") {
		leadsToPrefLen = len("tunnels lead to valves ")
	}
	leadsToValves := strings.Split(parts[1][leadsToPrefLen:], ", ")
	return ValveIDFromStr(valveIDStr),
		common.ParseInt[int](flowRateStr),
		common.Map(leadsToValves, ValveIDFromStr)
}

func parseValve(vm ValveMap, line string) {
	vId, flowRate, leadsTo := parseValveParams(line)
	_, exists := vm[vId]
	if !exists {
		vm[vId] = &Valve{
			ID: vId,
		}
	}
	vm[vId].FlowRate = flowRate
	for _, childValveID := range leadsTo {
		_, exists := vm[childValveID]
		if !exists {
			vm[childValveID] = &Valve{
				ID: childValveID,
			}
		}
		vm[vId].LeadsTo = append(vm[vId].LeadsTo, vm[childValveID])
	}
}

const minutesLimit = 30

func buildTurns(m ValveMap, from *Valve) <-chan Steps {
	ch := make(chan Steps)
	go func() {
		defer close(ch)

		allClosedValves := []*Valve{}
		for _, v := range m {
			if !v.IsOpen && v.FlowRate > 0 {
				allClosedValves = append(allClosedValves, v)
			}
		}
		for p := range common.PermutationsToChan(allClosedValves) {
			for _, v := range from.LeadsTo {
				if v == p[0] {
					ch <- buildSteps(p)
					break
				}
			}
		}
	}()
	return ch
}

func buildSteps(p []*Valve) Steps {
	var steps = Steps{{p[0], actionMove}, {p[0], actionOpen}}
	for i := 1; i < len(p); i++ {
		path := BfsPath(p[i-1], p[i])
		for j := range path {
			steps = append(steps, Step{path[j], actionMove})
		}
		steps = append(steps, Step{p[i], actionOpen})
	}
	return steps
}

func simulateTurn(m ValveMap, steps Steps) int {
	var totalPressure int
	for i := 0; i < minutesLimit; i++ {
		pres := m.TotalOpenedFlowRate()
		totalPressure += pres
		if i >= len(steps) {
			//fmt.Println(i, "wait", pres)
			continue
		}
		step := steps[i]
		//fmt.Println(i, step.String(), pres)
		if step.Action == actionOpen {
			step.Valve.IsOpen = true
		}
	}
	//fmt.Println("==============")
	return totalPressure
}

func simulate(m ValveMap, start *Valve) {
	bestScore := 0
	for turn := range buildTurns(m, start) {
		m.Reset()
		score := simulateTurn(m, turn)
		fmt.Println(score)
		if score > bestScore {
			bestScore = score
		}
	}
	fmt.Println("===")
	fmt.Println(bestScore)
}

func BfsPath(from, to *Valve) []*Valve {
	scores := map[ValveID]int{}
	scores[from.ID] = 0

	queue := []*Valve{from}
	visited := map[ValveID]bool{}
	parent := map[*Valve]*Valve{}

	var node *Valve
	for len(queue) > 0 {
		node, queue = queue[0], queue[1:]
		if visited[node.ID] {
			continue
		}
		visited[node.ID] = true
		if node == to {
			path := []*Valve{to}
			for path[0] != from {
				path = append([]*Valve{parent[path[0]]}, path...)
			}
			return path[1:]
		}
		for _, neighbor := range node.LeadsTo {
			if !visited[neighbor.ID] {
				scores[neighbor.ID] = scores[node.ID] + 1
				queue = append(queue, neighbor)
				parent[neighbor] = node
			}
		}
	}
	return nil
}

type action int8

const (
	actionOpen = action(iota)
	actionMove
)

func (a action) String() string {
	switch a {
	case actionMove:
		return "move to"
	case actionOpen:
		return "open"
	default:
		panic("bad action")
	}
}

type Step struct {
	Valve  *Valve
	Action action
}

type Steps []Step

func (s Steps) Score(currentMinute int) int {
	timeLeft := minutesLimit - currentMinute
	total := 0
	for i, step := range s {
		if step.Action == actionOpen {
			total += step.Valve.FlowRate * (timeLeft - i)
		}
	}
	return total
}

func (s Steps) Count() int {
	return len(s)
}

func (s Step) String() string {
	return fmt.Sprintf("%s %s", s.Action, s.Valve.ID)
}

func stepsToString(s Steps) string {
	return strings.Join(common.Map(s, func(step Step) string { return step.String() }), " - ")
}

func main() {
	sc, closeFile := common.FileScanner("./day16/input.txt")
	defer closeFile()
	var (
		valveMap = ValveMap{}
	)
	for sc.Scan() {
		parseValve(valveMap, sc.Text())
	}
	startValve := valveMap[ValveID{'A', 'A'}]
	simulate(valveMap, startValve)
}
