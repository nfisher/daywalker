package main

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
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
	Parent       *Dependency  `xml:"parent,omitempty"`
	Properties   Properties   `xml:"properties,omitempty"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
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

	for parent != nil {
		parentPom := seen[parent.Coord()]
		list = append(list, parentPom.Properties.List...)
		parent = parentPom.Parent
	}

	for _, props := range list {
		lookups["${"+props.Name()+"}"] = props.Value
	}

	c, err := p.deepCopy()
	if err != nil {
		return nil, err
	}

	for i, dep := range c.Dependencies {
		v, found := lookups[dep.Version]
		if found {
			dep.Version = v
			c.Dependencies[i] = dep
		}
	}

	return c, nil
}
