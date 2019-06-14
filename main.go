package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Position int

const (
	Head Position = iota
	Body
	Tail
)

func (p Position) String() string {
	switch p {
	case Head:
		return ""
	case Body:
		return " ├─ "
	case Tail:
		return " └─ "
	}
	panic("unknown position")
}

type Node struct {
	name     string
	children []*Node
}

func buildTree(records [][]string) *Node {
	if len(records) == 0 {
		return nil
	}
	linearMap := map[string][]string{}
	for _, record := range records {
		linearMap[record[0]] = append(linearMap[record[0]], record[1])
	}

	var f func(*Node)

	f = func(n *Node) {
		deps, ok := linearMap[n.name]
		if !ok {
			return
		}
		delete(linearMap, n.name)
		for _, dep := range deps {
			child := &Node{name: dep}
			f(child)
			n.children = append(n.children, child)
		}
	}

	curr := records[0][0]
	root := &Node{name: curr}
	f(root)

	return root
}

func PrintTree(out io.Writer, root *Node) {
	if root == nil {
		return
	}
	printTree(out, root, 0, Head, "")
}

func printTree(out io.Writer, root *Node, depth int, pos Position, prefix string) {
	fmt.Fprintf(out, "%s%s%s\n", prefix, pos.String(), root.name)
	if len(root.children) == 0 {
		return
	}

	var next string
	if pos == Tail {
		next = "    "
	} else if pos == Body {
		next = " │  "
	}
	newp := prefix + next

	cnt := len(root.children)
	for _, child := range root.children[:cnt-1] {
		printTree(out, child, depth+1, Body, newp)
	}
	printTree(out, root.children[cnt-1], depth+1, Tail, newp)
}

func main() {
	r := csv.NewReader(os.Stdin)
	r.Comma = ' '
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	tree := buildTree(records)
	PrintTree(os.Stdout, tree)
}
