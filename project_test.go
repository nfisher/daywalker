package main_test

import "testing"
import "github.com/nfisher/daywalker"

func Test_Relationship(t *testing.T) {
	var td = []struct {
		group    string
		artifact string
		version  string
		scope    string
		expected string
	}{
		{"ca.junctionbox", "wunderbar", "1234", "", "compile"},
		{"ca.junctionbox", "wunderbar", "1234", "compile", "compile"},
		{"ca.junctionbox", "wunderbar", "1234", "runtime", "runtime"},
		{"ca.junctionbox", "wunderbar", "1234", "test", "test"},
		{"ca.junctionbox", "wunderbar", "", "compile", "unresolved_compile"},
		{"ca.junctionbox", "wunderbar", "", "test", "unresolved_test"},
		{"ca.junctionbox", "wunderbar", "", "runtime", "unresolved_runtime"},
	}

	for _, data := range td {
		dep := &main.Dependency{
			GroupId:    data.group,
			ArtifactId: data.artifact,
			Version:    data.version,
			Scope:      data.scope,
		}
		actual := main.Relationship(dep)
		if data.expected != actual {
			t.Errorf("got Relationship(%v) = %v, want %v", dep, actual, data.expected)
		}
	}
}
