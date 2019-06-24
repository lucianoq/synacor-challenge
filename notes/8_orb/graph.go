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
//+----+---+----+----+
//|  * | 8 |  - | 1  |
//+------------------+
//|  4 | * | 11 | *  |
//+------------------+
//|  + | 4 |  - | 18 |
//+------------------+
//| 22 | - |  9 | *  |
//+----+---+----+----+
//
// IDs
//   6   7
// 4   5
//   2   3
// 0   1
func NewGraph() *Graph {
	g := &Graph{
		Nodes: make([]*Node, 0, 8),
		Edges: make([]*Edge, 0, 20),
	}

	g.Nodes = append(g.Nodes, &Node{0, 22})
	g.Nodes = append(g.Nodes, &Node{1, 9})
	g.Nodes = append(g.Nodes, &Node{2, 4})
	g.Nodes = append(g.Nodes, &Node{3, 18})
	g.Nodes = append(g.Nodes, &Node{4, 4})
	g.Nodes = append(g.Nodes, &Node{5, 11})
	g.Nodes = append(g.Nodes, &Node{6, 8})
	g.Nodes = append(g.Nodes, &Node{7, 1})

	g.Edges = append(g.Edges, &Edge{0, 1, Minus})
	g.Edges = append(g.Edges, &Edge{0, 2, Minus})
	g.Edges = append(g.Edges, &Edge{0, 2, Plus})
	g.Edges = append(g.Edges, &Edge{0, 4, Plus})

	g.Edges = append(g.Edges, &Edge{1, 0, Minus})
	g.Edges = append(g.Edges, &Edge{1, 2, Minus})
	g.Edges = append(g.Edges, &Edge{1, 5, Minus})
	g.Edges = append(g.Edges, &Edge{1, 3, Minus})
	g.Edges = append(g.Edges, &Edge{1, 3, Times})

	g.Edges = append(g.Edges, &Edge{2, 0, Plus})
	g.Edges = append(g.Edges, &Edge{2, 0, Minus})
	g.Edges = append(g.Edges, &Edge{2, 1, Minus})
	g.Edges = append(g.Edges, &Edge{2, 3, Minus})
	g.Edges = append(g.Edges, &Edge{2, 5, Minus})
	g.Edges = append(g.Edges, &Edge{2, 5, Times})
	g.Edges = append(g.Edges, &Edge{2, 4, Times})
	g.Edges = append(g.Edges, &Edge{2, 4, Plus})
	g.Edges = append(g.Edges, &Edge{2, 6, Times})

	g.Edges = append(g.Edges, &Edge{3, 1, Minus})
	g.Edges = append(g.Edges, &Edge{3, 1, Times})
	g.Edges = append(g.Edges, &Edge{3, 2, Minus})
	g.Edges = append(g.Edges, &Edge{3, 5, Minus})
	g.Edges = append(g.Edges, &Edge{3, 5, Times})
	g.Edges = append(g.Edges, &Edge{3, 7, Times})

	g.Edges = append(g.Edges, &Edge{4, 0, Plus})
	g.Edges = append(g.Edges, &Edge{4, 2, Plus})
	g.Edges = append(g.Edges, &Edge{4, 2, Times})
	g.Edges = append(g.Edges, &Edge{4, 6, Times})
	g.Edges = append(g.Edges, &Edge{4, 5, Times})

	g.Edges = append(g.Edges, &Edge{5, 2, Minus})
	g.Edges = append(g.Edges, &Edge{5, 2, Times})
	g.Edges = append(g.Edges, &Edge{5, 3, Minus})
	g.Edges = append(g.Edges, &Edge{5, 3, Times})
	g.Edges = append(g.Edges, &Edge{5, 6, Times})
	g.Edges = append(g.Edges, &Edge{5, 6, Minus})
	g.Edges = append(g.Edges, &Edge{5, 7, Minus})
	g.Edges = append(g.Edges, &Edge{5, 7, Times})
	g.Edges = append(g.Edges, &Edge{5, 4, Times})
	g.Edges = append(g.Edges, &Edge{5, 1, Minus})

	g.Edges = append(g.Edges, &Edge{6, 4, Times})
	g.Edges = append(g.Edges, &Edge{6, 5, Times})
	g.Edges = append(g.Edges, &Edge{6, 5, Minus})
	g.Edges = append(g.Edges, &Edge{6, 2, Times})
	g.Edges = append(g.Edges, &Edge{6, 7, Minus})

	return g
}
func (g *Graph) AddEdge(from, to int, op Operator) {
	g.Edges = append(g.Edges, &Edge{from, to, op})
	g.Edges = append(g.Edges, &Edge{6, 7, Minus})
}


func (g *Graph) walk(from, to int, visited [16]int, route string, depth int, acc int, maxDepth int) {
	if depth >= maxDepth {
		return
	}

	visited[from]++

	if from == to {
		if acc == 30 {
			fmt.Println(route)
			fmt.Println(acc)
			os.Exit(0)
		}
		return
	}

	for _, e := range g.Edges {
		if e.From == from {

			// We are not allowed to come back to the Antechamber
			if e.To != 0 {
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
}
