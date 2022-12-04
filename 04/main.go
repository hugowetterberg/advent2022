package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

type Range [2]int

func (r Range) Contains(b Range) bool {
	return r[0] <= b[0] && r[1] >= b[1]
}

func (r Range) Overlaps(b Range) bool {
	return r.Contains(b) || b.Contains(r) ||
		(r[0] >= b[0] && r[0] <= b[1]) ||
		(r[1] >= b[0] && r[1] <= b[1])
}

func run() error {
	var linum, containedPairs, overlappingPairs int

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := bytes.NewReader(r.Bytes())

		linum++

		var a, b Range

		_, err := fmt.Fscanf(line, "%d-%d,%d-%d",
			&a[0], &a[1], &b[0], &b[1])
		if err != nil {
			return fmt.Errorf("failed to parse line %d: %w",
				linum, err)
		}

		if a.Contains(b) || b.Contains(a) {
			containedPairs++
		}

		if a.Overlaps(b) {
			overlappingPairs++
		}
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	fmt.Printf("assignement pairs where one fully contains the other: %d\n",
		containedPairs)

	fmt.Printf("assignement pairs that overlap: %d\n",
		overlappingPairs)

	return nil
}
