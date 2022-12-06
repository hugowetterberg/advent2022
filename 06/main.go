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
	seq := make([]byte, 14)

	r := bufio.NewReader(os.Stdin)

	n, err := r.Read(seq)
	if err != nil {
		return fmt.Errorf("failed to read first four bytes: %w", err)
	}

	for {
		if uniq(seq) {
			break
		}

		copy(seq, seq[1:])

		b, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to read next byte: %w", err)
		}

		n++
		seq[len(seq)-1] = b
	}

	println("sequence is:", string(seq),
		"found after reading", n, "characters")

	return nil
}

func uniq(b []byte) bool {
	for i := 0; i < len(b); i++ {
		for j := i + 1; j < len(b); j++ {
			if b[i] == b[j] {
				return false
			}
		}
	}

	return true
}
