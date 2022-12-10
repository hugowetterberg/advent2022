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

const (
	hidden        = 0
	visibleTop    = 1
	visibleRight  = 1 << 1
	visibleBottom = 1 << 2
	visibleLeft   = 1 << 3
)

type Rule struct {
	Movement [2]int
	Mask     int
}

var rules = []Rule{
	{
		Movement: [2]int{0, 1},
		Mask:     visibleTop,
	},
	{
		Movement: [2]int{-1, 0},
		Mask:     visibleRight,
	},
	{
		Movement: [2]int{0, -1},
		Mask:     visibleBottom,
	},
	{
		Movement: [2]int{1, 0},
		Mask:     visibleLeft,
	},
}

func start(m [2]int, x, y int) [2]int {
	var s [2]int

	if m[0] < 0 {
		s[0] = 1
	}

	if m[1] < 0 {
		s[1] = 1
	}

	return [2]int{
		s[0] * x,
		s[1] * y,
	}
}

func next(s, i, m [2]int, x, y int) ([2]int, bool) {
	n := [2]int{
		i[0] + m[0],
		i[1] + m[1],
	}

	if m[0] != 0 && (n[0] < 0 || n[0] > x) {
		return shift(s, m, i), true
	}

	if m[1] != 0 && (n[1] < 0 || n[1] > y) {
		return shift(s, m, i), true
	}

	return n, false
}

func atEnd(m, i [2]int, x, y int) bool {
	if m[0] == 0 && (i[0] < 0 || i[0] > x) {
		return true
	}

	if m[1] == 0 && (i[1] < 0 || i[1] > y) {
		return true
	}

	return false
}

func shift(s [2]int, m [2]int, i [2]int) [2]int {
	if m[0] == 0 {
		return [2]int{i[0] + 1, s[1]}
	}

	return [2]int{s[0], i[1] + 1}
}

func run() error {
	var linum, colnum int

	treeMap := [][]int{}
	visibility := [][]int{}

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Bytes()
		linum++

		if colnum != 0 && len(line) != colnum {
			return fmt.Errorf("uneven grid on line %d", linum)
		}

		colnum = len(line)

		row := make([]int, colnum)

		for i := 0; i < colnum; i++ {
			row[i] = int(line[i] - 48)
		}

		treeMap = append(treeMap, row)

		vis := make([]int, colnum)
		visibility = append(visibility, vis)
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	for i := range rules {
		applyRule(rules[i], treeMap, visibility)
	}

	var visible int

	for i := 0; i < linum*colnum; i++ {
		y := i / colnum
		x := i % colnum

		if visibility[y][x] > 0 {
			visible++
		}
	}

	println("visible", visible)

	return nil
}

func applyRule(r Rule, m, v [][]int) {
	x := len(m[0]) - 1
	y := len(m) - 1
	s := start(r.Movement, x, y)
	height := -1

	var shifted bool

	for i := s; !atEnd(r.Movement, i, x, y); i, shifted = next(s, i, r.Movement, x, y) {
		current := m[i[1]][i[0]]

		if shifted {
			height = -1
		}

		if current > height {
			v[i[1]][i[0]] |= r.Mask
			height = current
		}
	}
}
