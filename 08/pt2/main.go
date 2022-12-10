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

func run() error {
	var linum, colnum int

	treeMap := [][]byte{}

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		if colnum != 0 && len(line) != colnum {
			return fmt.Errorf("uneven grid on line %d", linum)
		}

		colnum = len(line)

		row := make([]byte, colnum)

		for i := 0; i < colnum; i++ {
			row[i] = line[i] - 48
		}

		treeMap = append(treeMap, row)
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	var maxScore int

	for i := 0; i < linum*colnum; i++ {
		y := i / colnum
		x := i % colnum

		score := calculateScenicScore(x, y, treeMap)
		if score > maxScore {
			maxScore = score
		}
	}

	println("score", maxScore)

	return nil
}

func calculateScenicScore(x, y int, m [][]byte) int {
	var distances [4]int

	height := m[y][x]

	// Up
	for i := 1; y >= i; i++ {
		distances[0] = i

		if height <= m[y-i][x] {
			break
		}
	}

	// Down
	for i := 1; len(m) > y+i; i++ {
		distances[1] = i

		if height <= m[y+i][x] {
			break
		}
	}

	// Right
	for i := 1; len(m[y]) > x+i; i++ {
		distances[2] = i

		if height <= m[y][x+i] {
			break
		}
	}

	// left
	for i := 1; x >= i; i++ {
		distances[3] = i

		if height <= m[y][x-i] {
			break
		}
	}

	score := 1

	for _, d := range distances {
		score *= d
	}

	return score
}
