package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v\n", err)
		os.Exit(1)
	}
}

type Cell struct {
	Position Position
	Height   byte
	Cost     int
	Visited  bool
}

func (c *Cell) String() string {
	if c == nil {
		return "outside grid"
	}

	return fmt.Sprintf("%s %d %d",
		string([]byte{c.Height}),
		c.Position.X, c.Position.Y,
	)
}

type Grid [][]*Cell

func (g Grid) ToString(p Position) string {
	var b strings.Builder

	for _, row := range g {
		for _, cell := range row {
			symbol := string([]byte{cell.Height})
			if cell.Cost != 0 {
				symbol = strconv.Itoa(cell.Cost)
			}

			if cell.Visited {
				symbol = "."
			}

			if cell.Position == p {
				symbol = "X"
			}

			b.WriteString(symbol)
		}

		b.WriteRune('\n')
	}

	b.WriteRune('\n')

	return b.String()
}

func (g Grid) Get(p Position) *Cell {
	if p.X < 0 || p.Y < 0 || p.Y >= len(g) {
		return nil
	}

	row := g[p.Y]

	if p.X >= len(row) {
		return nil
	}

	return row[p.X]
}

type Position struct {
	X, Y int
}

func (p Position) Add(x, y int) Position {
	return Position{X: p.X + x, Y: p.Y + y}
}

func run() error {
	r := bufio.NewScanner(os.Stdin)

	var (
		evaluated   int
		grid        Grid
		candidates  []*Cell
		start, goal Position
	)

	for r.Scan() {
		line := r.Bytes()
		row := make([]*Cell, len(line))

		for i := range line {
			height := line[i]
			pos := Position{X: i, Y: len(grid)}

			switch height {
			case 'S':
				start = pos
				height = 'a'
			case 'E':
				goal = pos
				height = 'z'
			}

			c := Cell{
				Position: pos,
				Height:   height,
			}

			if c.Height == 'a' {
				c.Cost = 1
				candidates = append(candidates, &c)
			}

			row[i] = &c
		}

		grid = append(grid, row)
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	fmt.Printf("start: %s\n", grid.Get(start))
	fmt.Printf("goal: %s\n", grid.Get(goal))

	maxSteps := len(grid) * len(grid[0])
	position := start

	for evaluated <= maxSteps {
		candidates = append(candidates, getNeighbourCandidates(grid, position)...)

		if len(candidates) == 0 {
			return errors.New("no possible moves left")
		}

		sort.Slice(candidates, func(i, j int) bool {
			return (candidates[i].Height == 'a' && candidates[j].Height != 'a') ||
				candidates[i].Cost < candidates[j].Cost
		})

		current := candidates[0]

		println(grid.ToString(current.Position))

		position = current.Position
		current.Visited = true
		evaluated++

		copy(candidates, candidates[1:])
		candidates = candidates[0 : len(candidates)-1]

		time.Sleep(1 * time.Millisecond)

		if position == goal {
			break
		}

	}

	fmt.Printf("took %d steps\n", grid.Get(goal).Cost)

	return nil
}

func getNeighbourCandidates(grid Grid, pos Position) []*Cell {
	cell := grid.Get(pos)
	if cell == nil {
		return nil
	}

	check := []*Cell{
		grid.Get(pos.Add(0, -1)),
		grid.Get(pos.Add(1, 0)),
		grid.Get(pos.Add(0, 1)),
		grid.Get(pos.Add(-1, 0)),
	}

	var candidates []*Cell

	for _, c := range check {
		if c == nil || c.Cost != 0 {
			continue
		}

		if c.Height > cell.Height && c.Height-cell.Height > 1 {
			continue
		}

		if cell.Height == 'a' {
			c.Cost = 1
		} else {
			c.Cost = cell.Cost + 1
		}

		candidates = append(candidates, c)
	}

	return candidates
}
