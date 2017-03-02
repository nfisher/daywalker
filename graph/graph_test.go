package graph_test

import (
	"reflect"
	"testing"
)

import "github.com/nfisher/daywalker/graph"

func Test_Graph_Contains_Added_Node(t *testing.T) {
	g := graph.New()

	actual := g.Contains("b")
	if actual {
		t.Errorf("want g.Contains(a) = false, got %v", actual)
	}

	actual = g.Contains(graph.Root)
	if actual != true {
		t.Errorf("want g.Contains(Root) = true, got %v", actual)
	}
}

func Test_Graph_Len(t *testing.T) {
	g := graph.New()

	if g.Size() != 1 {
		t.Errorf("want g.Len() = 1, got %v", g.Size())
	}

	g.Edge(graph.Root, "a")

	if g.Size() != 2 {
		t.Errorf("want g.Len() = 2, got %v", g.Size())
	}
}

func Test_Graph_Children(t *testing.T) {
	g := graph.New()

	g.Edge(graph.Root, "a")
	g.Edge(graph.Root, "c")
	g.Edge("a", "b")

	var td = []struct {
		parent   string
		expected []string
	}{
		{graph.Root, []string{"a", "c"}},
		{"a", []string{"b"}},
		{"b", []string{}},
		{"c", []string{}},
	}

	for _, v := range td {
		actual := g.Children(v.parent)
		var names = make([]string, 0, 8)
		for _, n := range actual {
			names = append(names, n.Name())
		}

		if !reflect.DeepEqual(names, v.expected) {
			t.Errorf("want g.Children(%v) = %v, got %v", v.parent, v.expected, names)
		}
	}
}

func Test_Graph_Filtered_Children(t *testing.T) {
	g := graph.New()
	g.Edge(graph.Root, "b")
	g.Edge("b", "a", "parent")
	g.Edge("b", "c")

	b := g.Children("b", graph.HasRelationship("parent"))

	if len(b) != 1 {
		t.Errorf("want len(b) = 1, got %v", len(b))
	}
}

func Test_Graph_Edge_should_add_absent_nodes(t *testing.T) {
	g := graph.New()
	g.Edge("a", "b")

	if !g.Contains("a") {
		t.Errorf("want g.Contains(a) != false, got %v", g.Contains("a"))
	}

	if !g.Contains("b") {
		t.Errorf("want g.Contains(b) != false, got %v", g.Contains("b"))
	}
}
