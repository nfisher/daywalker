package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	coord := flag.Args()[0]

	seen := make(Set)
	FollowGraph(r, coord, seen)

	bf, err := os.Create("BUCK." + coord)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bf.Close()

	deps := make([]string, 0, len(seen))

	for k := range seen {
		err := r.RetrieveJar(k)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if k != coord {
			d := split(k)
			f := filepath.Base(filename(d, "jar"))
			line := fmt.Sprintf("prebuilt_jar(name='%v', binary_jar='%v', visibility=['PUBLIC'])\n\n", d.ArtifactId, f)
			deps = append(deps, d.ArtifactId)
			bf.Write([]byte(line))
		}
	}

	d := split(coord)
	f := filepath.Base(filename(d, "jar"))
	line := fmt.Sprintf("prebuilt_jar(name='%v', binary_jar='%v', visibility=['PUBLIC'], deps=[\n':%v'\n] )\n\n", d.ArtifactId, f, strings.Join(deps, "',\n':"))
	bf.Write([]byte(line))
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
