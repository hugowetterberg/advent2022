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

type RuneStack struct {
	top   int
	runes []rune
}

func (rs *RuneStack) Push(r ...rune) {
	rs.runes = append(rs.runes, r...)
	rs.top += len(r)
}

func (rs *RuneStack) Pop(n int) []rune {
	if rs.top == 0 || rs.top < n {
		return nil
	}

	rs.top -= n

	r := rs.runes[rs.top:]

	rs.runes = rs.runes[0:rs.top]

	return r
}

func (rs *RuneStack) Peek() rune {
	if rs.top == 0 {
		return 0
	}

	return rs.runes[rs.top-1]
}

func (rs *RuneStack) String() string {
	return string(rs.runes)
}

func run() error {
	var linum int

	var pileInput [][]rune

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		if len(line) == 0 {
			break
		}

		for i := 0; i*4 < len(line); i++ {
			if i == len(pileInput) {
				pileInput = append(pileInput, []rune{})
			}

			offset := i * 4

			pile := pileInput[i]

			if line[offset] == '[' {
				pile = append(pile, rune(line[offset+1]))
			}

			pileInput[i] = pile
		}
	}

	var runeStacks []*RuneStack

	for _, runes := range pileInput {
		var stack RuneStack

		for i := range runes {
			stack.Push(runes[len(runes)-1-i])
		}

		runeStacks = append(runeStacks, &stack)
	}

	printStacks(runeStacks)

	for r.Scan() {
		line := r.Bytes()
		linum++

		r := bytes.NewReader(line)

		var count, source, dst int

		_, err := fmt.Fscanf(r, "move %d from %d to %d",
			&count, &source, &dst)
		if err != nil {
			return fmt.Errorf("failed to parse line %d: %v",
				linum, err)
		}

		crates := runeStacks[source-1].Pop(count)
		if crates == nil {
			return fmt.Errorf("not enough crates in stack %d",
				source)
		}

		runeStacks[dst-1].Push(crates...)
	}

	err := r.Err()
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	println("Final arrangement")
	printStacks(runeStacks)

	for i, stack := range runeStacks {
		crate := stack.Peek()
		if crate == 0 {
			println(i, "is empty")
			continue
		}

		print(string(crate))
	}

	println()

	return nil
}

func printStacks(s []*RuneStack) {
	for i, stack := range s {
		fmt.Fprintf(os.Stdout, "%d: %s\n", i+1, stack.String())
	}
}
