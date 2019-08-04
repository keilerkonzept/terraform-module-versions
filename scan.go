package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
)

type scanner struct {
	Root    string
	Results []*moduleReference
	Paths   []string
}

const terraformSourceFileExt = ".tf"

func (s *scanner) ScanFile(path string) error {
	hclSource, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read terraform source %q: v", err)
	}
	var sourceFile struct{ Module map[string]*moduleReference }
	if err := hcl.Unmarshal(hclSource, &sourceFile); err != nil {
		return fmt.Errorf("process terraform source %q: v", err)
	}
	for k := range sourceFile.Module {
		m := sourceFile.Module[k]
		m.Path = path
		m.Name = k
		if err := m.ParseSource(); err != nil {
			log.Printf("parse module source: %v", err)
		}
		s.Results = append(s.Results, m)
		s.Paths = append(s.Paths, path)
	}
	return nil
}
