package main

import (
	"bufio"
	"fmt"

	"github.com/g-kap/advent-of-code-2022/common"
)

type Node struct {
	value      int
	prev, next *Node
	length     *int
}

func (n *Node) Eject() {
	n.next.prev, n.prev.next = n.prev, n.next
}

func (n *Node) InjectAfter(o *Node) {
	on := o.next
	o.next = n
	n.next = on
	n.prev = o
	n.next.prev = n
}

func (n *Node) Mix() {
	newPlace := n
	if n.value == 0 {
		return
	} else if n.value > 0 {
		newPlace = n.MoveForward(n.value % (*n.length - 1))
	} else if n.value < 0 {
		newPlace = n.MoveBack(-n.value % (*n.length - 1)).prev
	}
	//log.Printf("%d moves between %d and %d:\n", n.value, newPlace.value, newPlace.next.value)
	if n == newPlace {
		return
	}
	n.Eject()
	n.InjectAfter(newPlace)
	//log.Println(n.ForwardSlice())
}

func (n *Node) MoveForward(cnt int) *Node {
	nn := n
	for i := 0; i < cnt; i++ {
		nn = nn.next
	}
	return nn
}

func (n *Node) MoveBack(cnt int) *Node {
	nn := n
	for i := 0; i < cnt; i++ {
		nn = nn.prev
	}
	return nn
}

func (n *Node) ForwardSlice() []int {
	s := []int{}
	node := n
	for {
		s = append(s, node.value)
		node = node.next
		if node == n || node == nil {
			return s
		}
	}
}

func makeLinedList(sc *bufio.Scanner) []*Node {
	var nodes []*Node
	sc.Scan()
	length := 1
	firstNode := &Node{value: common.ParseInt[int](sc.Text()), length: &length}
	firstNode.next = firstNode
	firstNode.prev = firstNode
	lastNode := firstNode
	nodes = append(nodes, firstNode)
	for sc.Scan() {
		length++
		node := &Node{
			value:  common.ParseInt[int](sc.Text()),
			prev:   lastNode,
			length: &length,
		}
		nodes = append(nodes, node)
		lastNode.next = node
		lastNode = node
	}
	lastNode.next = firstNode
	firstNode.prev = lastNode
	return nodes
}

func detectCycle(start *Node) (bool, []*Node) {
	node := start
	visited := common.Set[*Node]{}
	path := []*Node{}
	var idx int
	var found bool
	for {
		path = append(path, node)
		if node.next == start {
			return false, nil
		}
		idx++
		if visited.Contains(node.next) {
			fmt.Println("cycle detected at node ")
			found = true
			break
		}
		visited.Add(node)
		if node.next == node {
			fmt.Println("cycle detected at node ")
			return true, path
		}
		if node.prev == node {
			fmt.Println("cycle detected at node ")
			found = true
			break
		}
		node = node.next
	}
	if !found {
		return false, nil
	}
	var (
		startIdx int
		sn       *Node
	)
	for startIdx, sn = range path {
		if sn == node.next {
			break
		}
	}
	return true, path[startIdx:]
}

func main() {
	sc, closeFile := common.DayFileScanner(20, false)
	defer closeFile()
	nodes := makeLinedList(sc)
	var zeroNode *Node

	for _, node := range nodes {
		if node.value == 0 {
			zeroNode = node
		}
		node.Mix()
	}

	n := zeroNode
	var ans []int
	for i := 0; i < 3; i++ {
		n = n.MoveForward(1000)
		ans = append(ans, n.value)
	}
	fmt.Println(ans, common.Sum(ans))
}
