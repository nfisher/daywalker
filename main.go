package main

import (
	"flag"
	"fmt"
)

type Set map[string]*Project

func (s Set) Contains(coord string) bool {
	_, ok := s[coord]
	return ok
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("need a starting coordinate (e.g. com.sparkjava:spark-core:2.5.4)")
		return
	}

	r := &Repository{
		Base: "https://search.maven.org/remotecontent?filepath=",
	}

	seen := make(Set)
	FollowGraph(r, flag.Args()[0], seen)

	for k, _ := range seen {
		err := r.RetrieveJar(k)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func FollowGraph(r *Repository, coord string, seen Set) {
	if seen.Contains(coord) {
		return
	}

	pom, err := r.RetrievePom(coord)
	if err != nil {
		fmt.Printf("%v - %v\n", coord, err)
		return
	}

	seen[coord] = pom

	if pom.Parent != nil {
		FollowGraph(r, pom.Parent.Coord(), seen)
	}

	merged, err := pom.MergeProperties(seen)
	if err != nil {
		fmt.Printf("%v - %v\n", coord, err)
		return
	}

	seen[coord] = merged

	for _, d := range merged.Dependencies {
		if d.Scope != "" {
			continue
		}

		if d.Version == "" {
			fmt.Println("SKIP - " + d.Coord())
			continue
		}

		FollowGraph(r, d.Coord(), seen)
	}
}
