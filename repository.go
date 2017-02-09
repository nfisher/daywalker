package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Repository struct {
	Base string
	sync.RWMutex
}

func split(coord string) *Dependency {
	s := strings.Split(coord, ":")
	return &Dependency{
		GroupId:    s[0],
		ArtifactId: s[1],
		Version:    s[2],
	}
}

func filename(d *Dependency, ext string) string {
	return fmt.Sprintf("%[1]v/%[3]v/%[2]v/%[3]v-%[2]v.%[4]v", strings.Replace(d.GroupId, ".", "/", -1), d.Version, d.ArtifactId, ext)
}

func (r *Repository) Url(asset string) string {
	r.RLock()
	base := r.Base
	r.RUnlock()

	return base + asset
}

func (r *Repository) RetrievePom(coord string) (*Project, error) {
	f := filename(split(coord), "pom")

	resp, err := http.Get(r.Url(f))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(resp.Body)
	project := &Project{}

	err = dec.Decode(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *Repository) RetrieveJar(coord string) error {
	d := split(coord)
	f := filename(d, "jar")

	resp, err := http.Get(r.Url(f))
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("want %v == 200, got %v\n", coord, resp.Status)
	}

	jf := fmt.Sprintf("%v-%v.jar", d.ArtifactId, d.Version)
	fmt.Println(jf)
	fd, err := os.Create(jf)
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
