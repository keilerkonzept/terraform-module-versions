package versions

import (
	"fmt"

	"github.com/sgreben/versions/pkg/semver"
)

type ConstraintForName struct {
	Name       string
	Constraint *semver.Constraints
}

type VersionForName struct {
	Name    string
	Version string
}

// ConstraintsForName maps a "package" name to its version constraint
type ConstraintsForName map[string]*semver.Constraints

// ConstraintsForVersion maps a package version to its set of dependency constraints
type ConstraintsForVersion map[string]ConstraintsForName

// ConstraintGraph is a map package->(version->(package->constraint))
type ConstraintGraph map[string]ConstraintsForVersion

// ConstraintStringGraph is a map package->(version->(package->constraint-string))
type ConstraintStringGraph map[string]map[string]map[string]string

func (g ConstraintGraph) Add(other ConstraintStringGraph) (err error) {
	for name, versionConstraints := range other {
		if g[name] == nil {
			g[name] = ConstraintsForVersion{}
		}
		for version, constraints := range versionConstraints {
			if g[name][version] == nil {
				g[name][version] = ConstraintsForName{}
			}
			for depName, constraint := range constraints {
				g[name][version][depName], err = semver.ParseConstraint(constraint)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// Solution maps a package name to a single selected version
type Solution struct {
	Selected map[string]*semver.Version
	Relaxed  map[string][]ConstraintForName
}

// SelectMVS uses the MVS (https://research.swtch.com/vgo-mvs) algorithm to solve a version constraint graph
func (d ConstraintGraph) SelectMVS() (*Solution, error) {
	var work []VersionForName
	for name, vcs := range d {
		if len(vcs) == 1 {
			for version := range vcs {
				work = append(work, VersionForName{
					Name:    name,
					Version: version,
				})
			}
		}
	}

	out := Solution{
		Selected: make(map[string]*semver.Version, len(d)),
		Relaxed:  map[string][]ConstraintForName{},
	}
	constraints := map[string][]ConstraintForName{}

	for len(work) > 0 {
		item := work[0]
		itemWants, ok := d[item.Name][item.Version]
		work = work[1:]

		if !ok {
			return nil, fmt.Errorf("no version %v of %q defined", item.Version, item.Name)
		}
		for name, wanted := range itemWants {
			var available semver.Collection
			for k := range d[name] {
				v, _ := semver.Parse(string(k))
				available = append(available, v)
			}
			constraints[name] = append(constraints[name], ConstraintForName{
				Name:       item.Name,
				Constraint: wanted,
			})
			matching := wanted.OldestMatching(available)
			if matching == nil {
				return nil, fmt.Errorf("no version of %q matching constraint %q of %q:%q (available: %v)", name, wanted, item.Name, item.Version, available)
			}
			if current, ok := out.Selected[name]; ok {
				if current.LessThan(matching) {
					out.Selected[name] = matching
					work = append(work, VersionForName{
						Name:    name,
						Version: matching.String(),
					})
				}
				continue
			}
			out.Selected[name] = matching
			work = append(work, VersionForName{
				Name:    name,
				Version: matching.String(),
			})
		}
	}

	for k, cs := range constraints {
		v := out.Selected[k]
		var es []ConstraintForName
		for _, c := range cs {
			if !c.Constraint.Check(v) {
				es = append(es, c)
			}
		}
		if len(es) == 0 {
			continue
		}
		for _, e := range es {
			out.Relaxed[e.Name] = append(out.Relaxed[e.Name], ConstraintForName{
				Name:       k,
				Constraint: e.Constraint,
			})
		}
	}
	return &out, nil
}
