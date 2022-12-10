package main

import (
	"bufio"
	"errors"
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

func (c Coord) FindAdjacent(b Coord) (Coord, error) {
	if c.Adjacent(b) {
		return c, nil
	}

	if c.X == b.X {
		c.Y += norm(c.Y, b.Y)
		return c, nil
	}

	if c.Y == b.Y {
		c.X += norm(c.X, b.X)
		return c, nil
	}

	for _, m := range diagonals {
		n := c.Add(m)

		if n.Adjacent(b) {
			return n, nil
		}
	}

	return Coord{}, errors.New("no possible move")
}

func (c Coord) Equals(b Coord) bool {
	return c.X == b.X && c.Y == b.Y
}

func norm(a, b int) int {
	if a > b {
		return -1
	}

	if a < b {
		return 1
	}

	return 0
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

var diagonals = map[string]Coord{
	"UR": {X: 1, Y: 1},
	"DR": {X: 1, Y: -1},
	"DL": {X: -1, Y: -1},
	"UL": {X: -1, Y: 1},
}

func run() error {
	var linum int

	vis := Visualiser{
		TopLeft:     Coord{Y: 4},
		BottomRight: Coord{X: 6},
	}

	flag.BoolVar(&vis.Enabled, "debug", false, "Enable debug visualiser")
	flag.Parse()

	knots := make([]Coord, 10)

	visitMap := map[string]int{}

	vis.Heading("Initial State")
	vis.Print(knots)

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
			knots[0] = knots[0].Add(m)

			for i := 1; i < len(knots); i++ {
				knots[i], err = knots[i].FindAdjacent(knots[i-1])
				if err != nil {
					return fmt.Errorf(
						"could not move knot %d: %w",
						i, err)
				}
			}

			vis.Print(knots)

			visitMap[knots[9].Key()]++
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

func (v Visualiser) Print(knots []Coord) {
	if !v.Enabled {
		return
	}

	for pos := v.TopLeft; pos.Y >= v.BottomRight.Y; pos.Y-- {
		for pos.X = v.TopLeft.X; pos.X < v.BottomRight.X; pos.X++ {
			printPos(pos, knots)
		}

		println()
	}

	println()
}

func printPos(p Coord, knots []Coord) {
	for i := 0; i < len(knots); i++ {
		if !p.Equals(knots[i]) {
			continue
		}

		if i == 0 {
			print("H")
			return
		}

		print(i)
		return
	}

	print(".")
}
