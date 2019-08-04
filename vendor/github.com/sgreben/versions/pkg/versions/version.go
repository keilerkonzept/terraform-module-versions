package versions

import (
	"sort"

	"github.com/sgreben/versions/pkg/semver"
)

// VersionWithSource is a semver version together with a source
type VersionWithSource struct {
	Version *semver.Version
	Source  VersionSource
}

// VersionSource is the source of a versioned artifact
type VersionSource struct {
	Git    *VersionSourceGit    `json:",omitempty"`
	Docker *VersionSourceDocker `json:",omitempty"`
}

// WithSources is a slice of VersionWithSource-s
type WithSources []VersionWithSource

// LatestMatching returns the latest version matching the given constraints
func (c WithSources) LatestMatching(constraints *semver.Constraints) *VersionWithSource {
	sort.Sort(c)
	for i := 0; i < len(c); i++ {
		candidate := c[len(c)-1-i]
		if constraints.Check(candidate.Version) {
			return &candidate
		}
	}
	return nil
}

// Len is sort.Interface.Len
func (c WithSources) Len() int {
	return len(c)
}

// Less is sort.Interface.Less
func (c WithSources) Less(i, j int) bool {
	return c[i].Version.LessThan(c[j].Version)
}

// Swap is sort.Interface.Swap
func (c WithSources) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
