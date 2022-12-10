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

var (
	cmdCd = "$ cd "
	cmdLs = "$ ls"
)

type Node struct {
	Name     string
	IsDir    bool
	Size     int
	Parent   *Node
	Children map[string]*Node
}

func NewNode(parent *Node) *Node {
	return &Node{
		Parent:   parent,
		Children: make(map[string]*Node),
	}
}

func run() error {
	var linum int

	root := NewNode(nil)
	root.IsDir = true

	cwd := root

	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Text()

		linum++

		if strings.HasPrefix(line, "$ ") {
			cmd, args, _ := strings.Cut(line[2:], " ")

			switch cmd {
			case "cd":
				if args == "" {
					return fmt.Errorf(
						"missing argument for cd on line %d",
						linum)
				}

				path := strings.Split(args, "/")

				for i, dir := range path {
					if i == 0 && dir == "" {
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
						child = NewNode(cwd)
						cwd.Children[dir] = child
					}

					cwd = child
				}
			case "ls":
				// We'll just ignore the actual command and
				// accept ls results to cwd as they come in
			default:
				return fmt.Errorf("unknown command %q on line %d",
					cmd, linum)

			}
		} else {
			ds, name, ok := strings.Cut(line, " ")
			if !ok {
				return fmt.Errorf("invalid dir entry %q on line %d",
					line, linum)
			}

			child := NewNode(cwd)
			child.Name = name

			if ds == "dir" {
				child.IsDir = true
			} else {
				size, err := strconv.Atoi(ds)
				if err != nil {
					return fmt.Errorf(
						"invalid file size %q on line %d",
						line, linum)
				}

				child.Size = size
			}

			cwd.Children[name] = child
		}
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	var sum int

	walkNodes(root, func(n *Node) {
		if n.Parent != nil {
			n.Parent.Size += n.Size
		}

		if n.IsDir && n.Size <= 100000 {
			sum += n.Size
		}
	})

	println("small directory sum", sum)

	var justRightSize int

	freeSpace := 70000000 - root.Size
	freeUp := 30000000 - freeSpace

	walkNodes(root, func(n *Node) {
		if !n.IsDir {
			return
		}

		if n.Size >= freeUp && (justRightSize == 0 || n.Size < justRightSize) {
			justRightSize = n.Size
		}
	})

	println("size of dir to delete", justRightSize)

	return nil
}

func walkNodes(n *Node, fn func(n *Node)) {
	for k := range n.Children {
		walkNodes(n.Children[k], fn)
	}

	fn(n)
}
