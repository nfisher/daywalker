package graph

import (
	"sync"
)

const Root = ""

type Node struct {
	name     string
	vertices []int
}

func (n *Node) Name() string {
	return n.name
}

func New() *Graph {
	g := &Graph{
		nodes: make([]Node, 0, 1024),
	}

	g.add(Root) // add root node.

	return g
}

type Graph struct {
	nodes []Node
	sync.RWMutex
}

func (g *Graph) Len() int {
	g.Lock()
	defer g.Unlock()
	return len(g.nodes)
}

func (g *Graph) Children(name string) ([]string, error) {
	g.Lock()
	defer g.Unlock()

	p := g.find(name)
	if p == -1 {
		return nil, nil
	}

	children := make([]string, 0, 8)
	parent := g.nodes[p]

	for _, i := range parent.vertices {
		children = append(children, g.nodes[i].Name())
	}

	return children, nil
}

func (g *Graph) Add(name string) {
	g.Lock()
	defer g.Unlock()

	g.add(name)
}

func (g *Graph) add(name string) int {
	i := g.find(name)
	if i != -1 {
		return i
	}

	n := Node{
		name:     name,
		vertices: make([]int, 0, 32),
	}

	g.nodes = append(g.nodes, n)

	return len(g.nodes) - 1
}

func (g *Graph) Edge(name1 string, name2 string) {
	g.Lock()
	defer g.Unlock()

	p := g.find(name1)
	c := g.find(name2)

	if p == -1 {
		p = g.add(name1)
	}

	if c == -1 {
		c = g.add(name2)
	}

	g.nodes[p].vertices = append(g.nodes[p].vertices, c)
}

func (g *Graph) Contains(name string) bool {
	g.RLock()
	defer g.RUnlock()

	return g.find(name) != -1
}

func (g *Graph) find(name string) int {
	for i, n := range g.nodes {
		if n.Name() == name {
			return i
		}
	}

	return -1
}
