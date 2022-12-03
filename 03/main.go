package main

import (
	"bufio"
	"errors"
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

func run() error {
	var linum, prioSum int

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		if len(line)%2 != 0 {
			return errors.New("a line cannot contain an odd number of items")
		}

		dupes := map[Item]bool{}

		for i := 0; i < len(line)/2; i++ {
			for j := len(line) / 2; j < len(line); j++ {
				if line[i] == line[j] {
					dupes[Item(line[i])] = true
				}
			}
		}

		for t := range dupes {
			prioSum += t.Priority()
		}
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	fmt.Printf("the sum of the priorities is %d\n", prioSum)

	return nil
}
