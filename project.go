package main

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"fmt"
)

type Dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

func (d *Dependency) Coord() string {
	return d.GroupId + ":" + d.ArtifactId + ":" + d.Version
}

type Property struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",chardata"`
}

func (p *Property) Name() string {
	return p.XMLName.Local
}

type Properties struct {
	List []Property `xml:",any"`
}

type Project struct {
	Parent               *Dependency  `xml:"parent,omitempty"`
	Properties           Properties   `xml:"properties,omitempty"`
	Dependencies         []Dependency `xml:"dependencies>dependency"`
	DependencyManagement []Dependency `xml:"dependencyManagement>dependencies>dependency"`
}

func (p *Project) deepCopy() (*Project, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}

	c := &Project{}
	err = dec.Decode(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *Project) MergeProperties(seen Set) (*Project, error) {
	lookups := make(map[string]string)
	list := p.Properties.List
	parent := p.Parent

	if parent != nil {
		lookups["${project.version}"] = parent.Version
	}

	deps := make([]Dependency, 0, 1024)
	depManagement := make([]Dependency, 0, 1024)
	coords := make(map[string]Dependency)

	for parent != nil {
		parentPom := seen[parent.Coord()]
		fmt.Printf(" RAWWWRR %v\n", parent.Coord())
		list = append(list, parentPom.Properties.List...)
		deps = append(deps, parentPom.Dependencies...)
		depManagement = append(depManagement, parentPom.DependencyManagement...)
		parent = parentPom.Parent
	}

	for _, c := range depManagement {
		// TODO (NF 2017-02-09): need to create a proper dependency graph.
		noVer := fmt.Sprintf("%v:%v:", c.GroupId, c.ArtifactId)
		coords[noVer] = c
	}

	for _, props := range list {
		lookups["${"+props.Name()+"}"] = props.Value
	}

	c, err := p.deepCopy()
	if err != nil {
		return nil, err
	}

	// import deps for BOM style projects
	c.Dependencies = append(c.Dependencies, deps...)

	for i, dep := range c.Dependencies {
		v, found := lookups[dep.Version]
		if found {
			dep.Version = v
			c.Dependencies[i] = dep
		}

		// if it's in the dependencyManagement section use that specification.
		d, found := coords[dep.Coord()]
		if found {
			c.Dependencies[i] = d
		}
	}

	return c, nil
}
