package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

type Move int

const (
	MoveRock    Move = 0
	MovePaper   Move = 1
	MoveScissor Move = 2
)

type Result int

const (
	Draw Result = 0
	Loss Result = -1
	Win  Result = 1
)

func (m Move) Score() int {
	return int(m) + 1
}

var resultMatrix = [][]Result{
	{Draw, Loss, Win},
	{Win, Draw, Loss},
	{Loss, Win, Draw},
}

var moveSymbolMap = map[byte]Move{
	'A': MoveRock,
	'B': MovePaper,
	'C': MoveScissor,
	'X': MoveRock,
	'Y': MovePaper,
	'Z': MoveScissor,
}

func (m Move) Compare(against Move) Result {
	return resultMatrix[m][against]
}

type Round struct {
	OpponentMove Move
	ResponseMove Move
}

func (r Round) Score() int {
	score := r.ResponseMove.Score()

	switch r.ResponseMove.Compare(r.OpponentMove) {
	case Win:
		score += 6
	case Draw:
		score += 3
	}

	return score
}

func run() error {
	var totalScore, linum int

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		var round Round

		if len(line) != 3 {
			return fmt.Errorf("invalid length %d (expected 3 )of line %d",
				len(line), linum)
		}

		opp, ok := moveSymbolMap[line[0]]
		if !ok {
			return fmt.Errorf("invalid opponent move %s",
				string(line[0]))
		}

		round.OpponentMove = opp

		resp, ok := moveSymbolMap[line[2]]
		if !ok {
			return fmt.Errorf("invalid response move %s",
				string(line[2]))
		}

		round.ResponseMove = resp

		totalScore += round.Score()
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	fmt.Printf("total score: %d\n", totalScore)

	return nil
}
