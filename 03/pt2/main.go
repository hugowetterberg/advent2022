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

type Item byte

func (it Item) IsValid() bool {
	return (it >= 65 && it <= 90) ||
		(it >= 97 && it <= 122)
}

func (it Item) Priority() int {
	if !it.IsValid() {
		return 0
	}

	if it < 91 {
		return int(it) - 38
	}

	return int(it) - 96
}

const presentInAll = 1 | 1<<1 | 1<<2

func run() error {
	var linum, prioSum int

	r := bufio.NewScanner(os.Stdin)

	var group [][]byte

	for r.Scan() {
		line := r.Bytes()
		linum++

		contents := make([]byte, len(line))
		copy(contents, line)

		group = append(group, contents)

		if len(group) < 3 {
			continue
		}

		priorities := make([]byte, 53)

		for i := range group {
			for _, b := range group[i] {
				priorities[Item(b).Priority()] |= 1 << i
			}
		}

		for prio, mask := range priorities {
			if mask == presentInAll {
				prioSum += prio
				break
			}
		}

		group = group[0:0]
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	fmt.Printf("the sum of the badges are %d\n", prioSum)

	return nil
}
