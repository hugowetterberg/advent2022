package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

func run() error {
	var packedCalories []int
	var elfSum, linum int

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		if len(line) == 0 {
			packedCalories = append(packedCalories, elfSum)
			elfSum = 0
			continue
		}

		calories, err := strconv.Atoi(string(line))
		if err != nil {
			return fmt.Errorf(
				"failed to parse calorie integer at line %d: %w",
				linum, err)
		}

		elfSum += calories
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	sort.Slice(packedCalories, func(i, j int) bool {
		return packedCalories[i] > packedCalories[j]
	})

	fmt.Printf("Max carried calories: %d\n", packedCalories[0])

	var topElvesSum int

	for i, elfSum := range packedCalories {
		if i == 3 {
			break
		}

		topElvesSum += elfSum

		fmt.Printf("%d. %d\n", i+1, elfSum)
	}

	fmt.Printf("The top three elves are carrying %d calories\n", topElvesSum)

	return nil
}
