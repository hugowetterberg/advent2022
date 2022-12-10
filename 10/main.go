package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

type Machine struct {
	X         int
	Operation Operation
	Failure   error
	Display   []bool
}

type Operation interface {
	Tick(m *Machine) bool
}

type NoOp struct{}

func (n NoOp) Tick(_ *Machine) bool {
	return true
}

type AddX struct {
	ticks int

	Number int
}

func (op *AddX) Tick(m *Machine) bool {
	op.ticks++

	if op.ticks < 2 {
		return false
	}

	m.X += op.Number

	return true
}

type Exception struct {
	Err error
}

func (op Exception) Tick(m *Machine) bool {
	m.Failure = op.Err

	return true
}

func run() error {
	r := bufio.NewScanner(os.Stdin)

	cancel := make(chan struct{})

	defer close(cancel)

	ops := make(chan Operation)

	go readInput(r, ops, cancel)

	var (
		rows = 6
		cols = 40
	)

	m := Machine{
		X:       1,
		Display: make([]bool, rows*cols),
	}

	var sum int

	for cycle := 1; true; cycle++ {
		if m.Operation == nil {
			in, ok := <-ops
			if ok {
				m.Operation = in
			}
		}

		if cycle == 20 || (cycle > 20 && (cycle-20)%40 == 0) {
			sum += cycle * m.X
		}

		position := (cycle - 1) % cols

		if m.X >= position-1 && m.X <= position+1 {
			m.Display[cycle-1] = true
		}

		if m.Operation == nil {
			break
		}

		if m.Operation.Tick(&m) {
			m.Operation = nil
		}

		if m.Failure != nil {
			return m.Failure
		}
	}

	println("sum", sum)

	for i, on := range m.Display {
		if on {
			print("#")
		} else {
			print(".")
		}

		if (i+1)%cols == 0 {
			println()
		}
	}

	return nil
}

func readInput(
	r *bufio.Scanner,
	instructions chan Operation,
	cancel chan struct{},
) {
	defer close(instructions)

	var linum int

	for r.Scan() {
		line := r.Text()
		linum++

		op, err := parseOperation(line)
		if err != nil {
			instructions <- Exception{
				Err: fmt.Errorf("parse operation %q on line %d: %w",
					line, linum, err),
			}

			return
		}

		select {
		case instructions <- op:
		case <-cancel:
			return
		}
	}

	if err := r.Err(); err != nil {
		instructions <- Exception{
			Err: fmt.Errorf("failed to read input: %w", err),
		}
	}
}

func parseOperation(line string) (Operation, error) {
	if line == "noop" {
		return NoOp{}, nil
	}

	op, arg, _ := strings.Cut(line, " ")

	switch op {
	case "addx":
		n, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid argument: %w", err)
		}

		return &AddX{
			Number: n,
		}, nil
	}

	return nil, fmt.Errorf("unknown operation %q", op)
}
