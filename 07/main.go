package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}

var (
	cmdCd = "$ cd "
	cmdLs = "$ ls"
)

type Node struct {
	IsDir bool
	Size int
	Parent *Node
	Children map[string]*Node
}

func NewNode(parent *Node) *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

func run() error {
	root := NewNode(nil)
	cwd := root
	
	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Text()

		if strings.HasPrefix(line, cmdCd) {			
				path := strings.Split(
					strings.TrimPrefix(line, cmdCd), "/")

			for i, dir := range path {
				if i== 0 && dir == "" {
					cwd = root
					continue
				}

				if dir == "." {
					continue
				}

				if dir == ".." {
					if cwd.Parent != nil {
						cwd = cwd.Parent
					}
					continue
				}

				child := cwd.Children[dir]
				if child == nil {
					child := NewNode(cwd)
				}
			}
		}
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	return nil
}

