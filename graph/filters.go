package graph

type Filter func(*Node, *Arc, *Node) bool

func Any() Filter {
	return func(from *Node, a *Arc, to *Node) bool {
		return true
	}
}

func HasRelationship(relationship string) Filter {
	return func(from *Node, a *Arc, to *Node) bool {
		for _, r := range a.relationships {
			if r == relationship {
				return true
			}
		}
		return false
	}
}
