package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/nfisher/daywalker/graph"
)

type Set map[string]*Project

func (s Set) Contains(coord string) bool {
	_, ok := s[coord]
	return ok
}

func main() {
	var coord string
	var useGraph bool

	flag.StringVar(&coord, "coord", "", "starting project coordinate [required]. (e.g. com.sparkjava:spark-core:2.5.4)")
	flag.BoolVar(&useGraph, "graph", false, "use digraph to map dependencies [in progress].")
	flag.Parse()

	if coord == "" {
		flag.Usage()
		return
	}

	r := &Repository{
		Base: "https://search.maven.org/remotecontent?filepath=",
	}

	if useGraph {
		var wg sync.WaitGroup

		ch := make(chan string, 32)
		g := graph.New()

		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for coord := range ch {
					fmt.Printf("[%v] start processing pom %v\n", i, coord)
					walkPoms(coord, r, g, ch)
					fmt.Printf("[%v] done processing pom %v \n", i, coord)
				}
			}(i)
		}

		ch <- coord

		wg.Wait()
		close(ch)

		graph.Print(g)
	} else {
		oldMain(coord, r)
	}
}

func processProperties(coord string, g *graph.Digraph, pom *Project) {
	for _, prop := range pom.Properties() {
		to := graph.NewNode("${" + prop.Name() + "}")
		to.Value = prop.Value()

		g.EdgeTo(coord, to, "property")
	}
}

func processParent(coord string, g *graph.Digraph, pom *Project, ch chan string) {
	if pom.Parent != nil {
		prop := graph.NewNode("${project.version}")
		prop.Value = pom.Parent.Version

		parentCoord := pom.Parent.Coord()

		g.Edge(coord, parentCoord, "parent")
		g.Edge(parentCoord, coord, "child")
		g.EdgeTo(coord, prop, "property")

		ch <- parentCoord
	}
}

func processManagedDependencies(coord string, g *graph.Digraph, pom *Project) {
	for _, dm := range pom.DependencyManagement {
		depCoord := dm.Coord()

		rel := ManagedRelationship(dm)

		g.Edge(coord, depCoord, rel)
	}
}

func processDependencies(coord string, g *graph.Digraph, pom *Project) {
	for _, dep := range pom.Dependencies {
		depCoord := dep.Coord()

		if strings.Contains(depCoord, "$") {
			properties := g.Children(coord, graph.HasRelationship("property"))
			for _, p := range properties {
				if strings.Contains(depCoord, p.Name()) {
					v, ok := p.Value.(string)
					if !ok {
						fmt.Printf("property %#v has unexpected value %#v\n", p.Name(), p.Value)
					}
					depCoord = strings.Replace(depCoord, p.Name(), v, -1)
					dep = split(depCoord)
				}
			}
		}

		rel := Relationship(dep)

		g.Edge(coord, depCoord, rel)
	}
}

var seen map[string]struct{} = make(map[string]struct{})

type SeenSet struct {
	seen map[string]struct{}
	sync.Mutex
}

func (ss *SeenSet) Add(v string) {
	ss.seen[v] = struct{}{}
}

func (ss *SeenSet) Contains(v string) bool {
	_, ok := ss.seen[v]

	return ok
}

var walked *SeenSet = &SeenSet{seen: make(map[string]struct{})}

func walkPoms(coord string, r *Repository, g *graph.Digraph, ch chan string) {
	walked.Lock()
	if walked.Contains(coord) {
		return
	}

	walked.Add(coord)
	walked.Unlock()

	pom, err := r.RetrievePom(coord)
	if err != nil {
		log.Println(err)
		return
	}

	processProperties(coord, g, pom)

	processParent(coord, g, pom, ch)

	processManagedDependencies(coord, g, pom)

	processDependencies(coord, g, pom)

	graph.Print(g)

	for _, depNode := range g.Children(coord, graph.HasRelationship("compile")) {
		ch <- depNode.Name()
	}
}
