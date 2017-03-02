package graph

func NewNode(name string) *Node {
	return &Node{
		name: name,
		arcs: make([]*Arc, 0, 8),
	}
}

type Node struct {
	name  string
	arcs  []*Arc
	Value interface{}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Arcs() []*Arc {
	return n.arcs
}

func (n *Node) AddArc(to int, relationships ...string) {
	n.arcs = append(n.arcs, &Arc{to, relationships})
}

type Arc struct {
	nodeIndex     int
	relationships []string
}
