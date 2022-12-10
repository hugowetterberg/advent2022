package main

import (
	"bufio"
	"flag"
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

type Coord struct {
	X int
	Y int
}

func (c Coord) Add(v Coord) Coord {
	return Coord{
		X: c.X + v.X,
		Y: c.Y + v.Y,
	}
}

func (c Coord) Key() string {
	return fmt.Sprintf("%d:%d", c.X, c.Y)
}

func (c Coord) Adjacent(b Coord) bool {
	return delta(c.X, b.X) <= 1 && delta(c.Y, b.Y) <= 1
}

func (c Coord) Equals(b Coord) bool {
	return c.X == b.X && c.Y == b.Y
}

func delta(a, b int) int {
	if a > b {
		return a - b
	}

	return b - a
}

var movements = map[string]Coord{
	"U": {X: 0, Y: 1},
	"R": {X: 1, Y: 0},
	"D": {X: 0, Y: -1},
	"L": {X: -1, Y: 0},
}

func run() error {
	var linum int

	vis := Visualiser{
		TopLeft:     Coord{Y: 4},
		BottomRight: Coord{X: 6},
	}

	flag.BoolVar(&vis.Enabled, "debug", false, "Enable debug visualiser")
	flag.Parse()

	var head, prevHead, tail Coord

	visitMap := map[string]int{}

	vis.Heading("Initial State")
	vis.Print(head, tail)

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Text()
		linum++

		vis.Heading(line)

		move, distance, ok := strings.Cut(line, " ")
		if !ok {
			return fmt.Errorf("invalid move %q on line %d",
				line, linum)
		}

		d, err := strconv.Atoi(distance)
		if err != nil {
			return fmt.Errorf(
				"invalid move distance for %q on line %d: %w",
				line, linum, err)

		}

		m, ok := movements[move]
		if !ok {
			return fmt.Errorf(
				"invalid move type for %q on line %d",
				line, linum)
		}

		for i := 0; i < d; i++ {
			prevHead = head
			head = head.Add(m)

			if !tail.Adjacent(head) {
				tail = prevHead
			}

			vis.Print(head, tail)

			visitMap[tail.Key()]++
		}
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	visited := len(visitMap)

	println("visited", visited)

	return nil
}

type Visualiser struct {
	Enabled     bool
	TopLeft     Coord
	BottomRight Coord
}

func (v Visualiser) Heading(txt string) {
	if !v.Enabled {
		return
	}

	println("==", txt, "==\n")
}

func (v Visualiser) Print(head, tail Coord) {
	if !v.Enabled {
		return
	}

	for pos := v.TopLeft; pos.Y >= v.BottomRight.Y; pos.Y-- {
		for pos.X = v.TopLeft.X; pos.X < v.BottomRight.X; pos.X++ {
			if pos.Equals(head) {
				print("H")
			} else if pos.Equals(tail) {
				print("T")
			} else {
				print(".")
			}
		}

		println()
	}

	println()
}
