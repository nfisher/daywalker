package graph

import (
	"fmt"
	"sync"
)

const Root = ""
const NotFound = -1

func New() *Digraph {
	g := &Digraph{
		nodes: make([]*Node, 0, 1024),
	}

	g.add(Root) // add root node.

	return g
}

func Print(g *Digraph) {
	g.Lock()
	defer g.Unlock()
	for _, n := range g.nodes {
		fmt.Printf("%v\n", n.Name())
		for j, arc := range n.Arcs() {
			i := arc.nodeIndex
			to := g.nodes[i]
			if len(n.Arcs()) == j+1 {
				fmt.Printf("  +-- %v -> %v\n", arc.relationships, to.Name())
			} else {
				fmt.Printf("  |-- %v -> %v\n", arc.relationships, to.Name())
			}
		}
	}
}

type Digraph struct {
	nodes []*Node
	sync.RWMutex
}

func (g *Digraph) Find(name string) *Node {
	g.RLock()
	defer g.RUnlock()

	i := g.find(name)

	if i == NotFound {
		return nil
	}

	return g.get(i)
}

func (g *Digraph) Size() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.nodes)
}

func (g *Digraph) Children(name string, filters ...Filter) []*Node {
	g.RLock()
	defer g.RUnlock()

	p := g.find(name)
	if p == NotFound {
		return nil
	}

	children := make([]*Node, 0, 8)
	parent := g.get(p)

	if len(filters) < 1 {
		filters = append(filters, Any())
	}

	for _, a := range parent.Arcs() {
		n := g.get(a.nodeIndex)
		if n == nil {
			fmt.Println("skipping nil node")
			continue
		}
		for _, f := range filters {
			if f(parent, a, n) {
				children = append(children, n)
			}
		}
	}

	return children
}

func (g *Digraph) get(i int) *Node {
	if i >= len(g.nodes) {
		fmt.Printf("wtf %v\n", i)
		return nil
	}
	return g.nodes[i]
}

func (g *Digraph) Add(name string) {
	g.Lock()
	defer g.Unlock()

	g.add(name)
}

func (g *Digraph) find(name string) int {
	for i, n := range g.nodes {
		if n.Name() == name {
			return i
		}
	}

	return NotFound
}

func (g *Digraph) add(name string) int {
	i := g.find(name)
	if i != NotFound {
		return i
	}

	n := NewNode(name)

	g.nodes = append(g.nodes, n)

	return len(g.nodes) - 1
}

func (g *Digraph) EdgeTo(f string, to *Node, relationships ...string) {
	if to == nil {
		return
	}

	g.Lock()
	defer g.Unlock()

	from := g.find(f)

	if from == NotFound {
		from = g.add(f)
	}

	pos := len(g.nodes)
	g.nodes = append(g.nodes, to)

	g.nodes[from].AddArc(pos, relationships...)
}

func (g *Digraph) Edge(f string, t string, relationships ...string) {
	g.Lock()
	defer g.Unlock()

	from := g.find(f)
	to := g.find(t)

	if from == NotFound {
		from = g.add(f)
	}

	if to == NotFound {
		to = g.add(t)
	}

	g.nodes[from].AddArc(to, relationships...)
}

func (g *Digraph) Contains(name string) bool {
	g.RLock()
	defer g.RUnlock()

	return g.find(name) != -1
}
