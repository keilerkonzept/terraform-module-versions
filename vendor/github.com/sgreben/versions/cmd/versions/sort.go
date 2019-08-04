package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/sgreben/versions/pkg/semver"
)

func sortCmd(versions []string, limit int) {
	svs := make(semver.Collection, 0, len(versions))
	for _, v := range versions {
		sv, err := semver.Parse(v)
		if err != nil {
			exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf(`"%s": %v`, v, err))
			continue
		}
		svs = append(svs, sv)
	}
	sort.Sort(svs)
	if limit > 0 && len(svs) > limit {
		svs = svs[len(svs)-limit:]
	}
	err := jsonEncode(svs, os.Stdout)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
	}
}
