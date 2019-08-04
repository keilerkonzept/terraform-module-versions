package main

import (
	"fmt"
	"os"

	"github.com/sgreben/versions/pkg/semver"
	"github.com/sgreben/versions/pkg/versions"
)

func selectSingleCmd(constraint string, versions []string) {
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("cannot parse constraint %q: %v", constraint, err))
		return
	}
	svs := make(semver.Collection, 0, len(versions))
	for _, v := range versions {
		sv, err := semver.Parse(v)
		if err != nil {
			exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf(`"%s": %v`, v, err))
			continue
		}
		svs = append(svs, sv)
	}
	solution := c.LatestMatching(svs)
	if solution == nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, "no matching version")
		return
	}
	jsonEncode(solution.String(), os.Stdout)
}

func selectAllCmd(constraint string, versions []string) {
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("cannot parse constraint %q: %v", constraint, err))
		return
	}
	svs := make(semver.Collection, 0, len(versions))
	for _, v := range versions {
		sv, err := semver.Parse(v)
		if err != nil {
			exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf(`"%s": %v`, v, err))
			continue
		}
		svs = append(svs, sv)
	}
	solution := c.AllMatching(svs)
	jsonEncode(solution, os.Stdout)
}

func selectMvsCmd(g versions.ConstraintGraph) {
	mvsSolution, err := g.SelectMVS()
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("mvs failed: %v", err))
		return
	}
	var out struct {
		Selected   map[string]string
		Relaxed map[string]map[string]string
	}
	out.Selected = make(map[string]string, len(mvsSolution.Selected))
	out.Relaxed = make(map[string]map[string]string, len(mvsSolution.Relaxed))
	for k, v := range mvsSolution.Selected {
		out.Selected[k] = v.String()
	}
	for k, es := range mvsSolution.Relaxed {
		m := make(map[string]string, len(es))
		out.Relaxed[k] = m
		for _, e := range es {
			m[e.Name] = e.Constraint.String()
		}
	}
	jsonEncode(out, os.Stdout)
}
