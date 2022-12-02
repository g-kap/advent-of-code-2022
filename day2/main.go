package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Sign int

const (
	Stone    = Sign(1)
	Paper    = Sign(2)
	Scissors = Sign(3)
)

var winCombinations = [9][2]Sign{
	{Stone, Scissors},
	{Paper, Stone},
	{Scissors, Paper},
}

type SignCombination struct {
	PlayerSign   Sign
	OpponentSign Sign
}

func (sc SignCombination) Score() (int, int) {
	var pScore, opScore int
	if sc.PlayerSign == sc.OpponentSign {
		pScore, opScore = 3, 3
	} else {
		for _, winComb := range winCombinations {
			if [2]Sign{sc.PlayerSign, sc.OpponentSign} == winComb {
				pScore, opScore = 6, 0
				break
			}
			if [2]Sign{sc.OpponentSign, sc.PlayerSign} == winComb {
				pScore, opScore = 0, 6
				break
			}
		}
	}
	return pScore + int(sc.PlayerSign), opScore + int(sc.OpponentSign)
}

func parseCombination(txt string, part2 bool) (SignCombination, error) {
	txt = strings.Replace(txt, " ", "", -1)
	var comb SignCombination
	if len(txt) != 2 {
		return comb, errors.New("could not parse signs combination")
	}
	switch txt[0] {
	case 'A':
		comb.OpponentSign = Stone
	case 'B':
		comb.OpponentSign = Paper
	case 'C':
		comb.OpponentSign = Scissors
	default:
		return comb, errors.New("could not parse signs combination")
	}
	if part2 {
		var err error
		comb.PlayerSign, err = selectSign(comb.OpponentSign, CombatResult(txt[1]))
		if err != nil {
			return comb, err
		}
	} else {
		switch txt[1] {
		case 'X':
			comb.PlayerSign = Stone
		case 'Y':
			comb.PlayerSign = Paper
		case 'Z':
			comb.PlayerSign = Scissors
		default:
			return comb, errors.New("could not parse signs combination")
		}
	}
	return comb, nil
}

type CombatResult byte

const (
	Loose = CombatResult('X')
	Draw  = CombatResult('Y')
	Win   = CombatResult('Z')
)

func selectSign(opSign Sign, result CombatResult) (Sign, error) {
	if result == Draw {
		return opSign, nil
	}
	for _, comb := range winCombinations {
		if result == Loose && comb[0] == opSign {
			return comb[1], nil
		}
		if result == Win && comb[1] == opSign {
			return comb[0], nil
		}
	}
	return Sign(0), errors.New("unsupported sign was provided")
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("can not open ./input.txt")
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var playerTotal int
	for sc.Scan() {
		comb, err := parseCombination(sc.Text(), true)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		pScore, _ := comb.Score()
		playerTotal += pScore
	}
	fmt.Println(playerTotal)
}
