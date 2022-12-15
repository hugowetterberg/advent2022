package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

type Monkey struct {
	ID           int
	Items        []int
	Operation    func(n int) int `json:"-"`
	TestDiv      int
	TrueMonkey   int
	FalseMonkey  int
	InspectCount int
}

type prefixParse struct {
	Name   string
	Prefix string
	Parse  func(line string, m *Monkey) error
}

var monkeyParsers = []prefixParse{
	{
		Name:   "starting items",
		Prefix: "  Starting items: ",
		Parse:  parseStarters,
	},
	{
		Name:   "operation",
		Prefix: "  Operation: ",
		Parse:  parseOperation,
	},
	{
		Name:   "test",
		Prefix: "  Test: ",
		Parse: func(line string, m *Monkey) error {
			return parseNum(
				line, "divisible by %d", &m.TestDiv)
		},
	},
	{
		Name:   "true action",
		Prefix: "    If true: ",
		Parse: func(line string, m *Monkey) error {
			return parseNum(
				line, "throw to monkey %d", &m.TrueMonkey)
		},
	},
	{
		Name:   "false action",
		Prefix: "    If false: ",
		Parse: func(line string, m *Monkey) error {
			return parseNum(
				line, "throw to monkey %d", &m.FalseMonkey)

		},
	},
}

func run() error {
	r := bufio.NewScanner(os.Stdin)

	var monkeys []*Monkey

	for r.Scan() {
		monkeyDecl := r.Text()

		var monkey Monkey

		_, err := fmt.Sscanf(monkeyDecl, "Monkey %d:", &monkey.ID)
		if err != nil {
			return fmt.Errorf(
				"failed to parse monkey declaraton %q: %w",
				monkeyDecl, err)
		}

		for _, p := range monkeyParsers {
			if !r.Scan() {
				return fmt.Errorf(
					"failed to read line for %s: %w",
					p.Name, r.Err())
			}

			line := r.Text()

			if !strings.HasPrefix(line, p.Prefix) {
				return fmt.Errorf(
					"expected %s line to start with %q, got %q",
					p.Name, p.Prefix, line)
			}

			value := strings.TrimPrefix(line, p.Prefix)

			err := p.Parse(value, &monkey)
			if err != nil {
				return fmt.Errorf("failed to parse %s value: %w",
					p.Name, err)
			}
		}

		monkeys = append(monkeys, &monkey)

		// Consume empty line
		r.Scan()
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	for round := 0; round < 20; round++ {
		for _, m := range monkeys {
			m.InspectCount += len(m.Items)

			for _, n := range m.Items {
				n = m.Operation(n) / 3

				var recipientIdx int

				if n%m.TestDiv == 0 {
					recipientIdx = m.TrueMonkey
				} else {
					recipientIdx = m.FalseMonkey
				}

				recipient := monkeys[recipientIdx]

				recipient.Items = append(recipient.Items, n)
			}

			m.Items = m.Items[0:0]
		}
	}

	for i := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n",
			i, monkeys[i].InspectCount)
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectCount > monkeys[j].InspectCount
	})

	fmt.Printf("Level of monkey business: %d\n",
		monkeys[0].InspectCount*monkeys[1].InspectCount)

	return nil
}

func parseStarters(line string, m *Monkey) error {
	s := strings.Split(line, ", ")
	items := make([]int, len(s))

	for i := range s {
		n, err := strconv.Atoi(s[i])
		if err != nil {
			return fmt.Errorf(
				"failed to parse starting item %d: %w",
				i+1, err)
		}

		items[i] = n
	}

	m.Items = items

	return nil
}

func parseNum(line, format string, n *int) error {
	_, err := fmt.Sscanf(line, format, n)
	if err != nil {
		return fmt.Errorf("failed to action %q as %q", line, format)
	}

	return nil
}

func parseOperation(line string, m *Monkey) error {
	var (
		op, vRef string
		opFunc   func(a, b int) int
	)

	const format = "new = old %s %s"

	_, err := fmt.Sscanf(line, format, &op, &vRef)
	if err != nil {
		return fmt.Errorf("failed to parse %q as %q", line, format)
	}

	switch op {
	case "*":
		opFunc = func(a, b int) int {
			return a * b
		}
	case "+":
		opFunc = func(a, b int) int {
			return a + b
		}
	default:
		return fmt.Errorf("unknown operator %q", op)
	}

	switch vRef {
	case "old":
		m.Operation = func(old int) int {
			v := opFunc(old, old)
			return v
		}
	default:
		v, err := strconv.Atoi(vRef)
		if err != nil {
			return fmt.Errorf("failed to parse %q as an integer: %w",
				vRef, err)
		}

		m.Operation = func(old int) int {
			return opFunc(old, v)
		}
	}

	return nil
}
