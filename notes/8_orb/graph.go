package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Operator string

const (
	Plus  Operator = "+"
	Minus Operator = "-"
	Times Operator = "*"
)

type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

type Node struct {
	ID    int
	Value int
}

type Edge struct {
	From  int
	To    int
	Label Operator
}

func main() {
	g := NewGraph()

	// Try different max depth, increasing it
	// Program will exit at the first solution
	for i := 3; i < 100; i++ {
		g.walk(0, 7, [16]int{}, "22", 0, 22, i)
	}
}

// Values
//   +----+----+----+----+
//   |  * |  8 |  - |  1 |
//   +----+----+----+----+
//   |  4 |  * | 11 |  * |
//   +----+----+----+----+
//   |  + |  4 |  - | 18 |
//   +----+----+----+----+
//   | 22 |  - |  9 |  * |
//   +----+----+----+----+
//
// IDs
//   +---+---+---+---+
//   |   | 6 |   | 7 |
//   +---+---+---+---+
//   | 4 |   | 5 |   |
//   +---+---+---+---+
//   |   | 2 |   | 3 |
//   +---+---+---+---+
//   | 0 |   | 1 |   |
//   +---+---+---+---+
func NewGraph() *Graph {
	g := &Graph{
		Nodes: make([]*Node, 0, 8),
		Edges: make([]*Edge, 0, 24),
	}

	g.Nodes = []*Node{{0, 22}, {1, 9}, {2, 4}, {3, 18}, {4, 4}, {5, 11}, {6, 8}, {7, 1}}

	g.AddEdge(0, 1, Minus)
	g.AddEdge(0, 2, Minus)
	g.AddEdge(0, 2, Plus)
	g.AddEdge(0, 4, Plus)
	g.AddEdge(1, 2, Minus)
	g.AddEdge(1, 3, Minus)
	g.AddEdge(1, 3, Times)
	g.AddEdge(1, 5, Minus)
	g.AddEdge(2, 3, Minus)
	g.AddEdge(2, 4, Times)
	g.AddEdge(2, 4, Plus)
	g.AddEdge(2, 5, Minus)
	g.AddEdge(2, 5, Times)
	g.AddEdge(2, 6, Times)
	g.AddEdge(3, 5, Minus)
	g.AddEdge(3, 5, Times)
	g.AddEdge(3, 7, Times)
	g.AddEdge(4, 5, Times)
	g.AddEdge(4, 6, Times)
	g.AddEdge(5, 6, Times)
	g.AddEdge(5, 6, Minus)
	g.AddEdge(5, 7, Minus)
	g.AddEdge(5, 7, Times)
	g.AddEdge(6, 7, Minus)

	return g
}
func (g *Graph) AddEdge(from, to int, op Operator) {
	g.Edges = append(g.Edges, &Edge{from, to, op})

	// We are not allowed to come back to the Antechamber
	// Or to leave the vault door room 
	if from != 0 && to != 7 {
		g.Edges = append(g.Edges, &Edge{to, from, Minus})
	}
}

func (g *Graph) walk(from, to int, visited [16]int, route string, depth int, acc int, maxDepth int) {
	if depth >= maxDepth {
		return
	}

	visited[from]++

	if from == to {
		if acc == 30 {
			fmt.Println(route)
			os.Exit(0)
		}
		return
	}

	for _, e := range g.Edges {
		if e.From == from {

			path := route + string(e.Label) + strconv.Itoa(g.Nodes[e.To].Value)

			var accumulator int
			switch e.Label {
			case Plus:
				accumulator = acc + g.Nodes[e.To].Value
			case Minus:
				accumulator = acc - g.Nodes[e.To].Value
			case Times:
				accumulator = acc * g.Nodes[e.To].Value
			default:
				log.Fatal()
			}

			g.walk(e.To, to, visited, path, depth+1, accumulator, maxDepth)
		}
	}
}
