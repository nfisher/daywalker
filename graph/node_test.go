package graph_test

import (
	"reflect"
	"testing"
)
import "github.com/nfisher/daywalker/graph"

func Test_Graph_Contains_Added_Node(t *testing.T) {
	g := graph.New()

	if g.Contains("b") {
		t.Errorf("want g.Contains(a) = false, got %v", g.Contains("b"))
	}

	g.Add("a")
	if !g.Contains("a") {
		t.Errorf("want g.Contains(a) = true, got %v", g.Contains("a"))
	}
}

func Test_Graph_Len(t *testing.T) {
	g := graph.New()

	if g.Len() != 1 {
		t.Errorf("want g.Len() = 1, got %v", g.Len())
	}

	g.Edge(graph.Root, "a")

	if g.Len() != 2 {
		t.Errorf("want g.Len() = 2, got %v", g.Len())
	}
}

func Test_Graph_Children(t *testing.T) {
	g := graph.New()

	expected := []string{"a", "c"}
	g.Edge(graph.Root, "a")
	g.Edge(graph.Root, "c")
	actual, _ := g.Children(graph.Root)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("want g.Children(Root) = [a], got %v", actual)
	}

	expected = []string{"b"}
	g.Edge("a", "b")
	actual, _ = g.Children("a")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("want g.Children(a) = [b], got %v", actual)
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
