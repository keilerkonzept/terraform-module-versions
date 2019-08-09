package versions

import (
	"github.com/sgreben/versions/pkg/semver"
	"github.com/sgreben/versions/pkg/simpledocker"
	"github.com/sgreben/versions/pkg/simplegit"
)

// Source is a source of versions
type Source struct {
	Git    *SourceGit
	Docker *SourceDocker
}

// Fetch retrieves versions
func (t Source) Fetch() (WithSources, error) {
	switch {
	case t.Git != nil:
		return t.Git.Fetch()
	case t.Docker != nil:
		return t.Docker.Fetch()
	}
	return nil, nil
}

// SourceGit is a git repository used as a version source
type SourceGit struct {
	Repository simplegit.Repository
}

// Fetch retrieves git tags as versions
func (t SourceGit) Fetch() (WithSources, error) {
	tags, err := t.Repository.RemoteRefs()
	if err != nil {
		return nil, err
	}
	out := make(WithSources, 0, len(tags))
	for _, tag := range tags {
		version, err := semver.Parse(tag.Name)
		if err != nil {
			continue
		}
		source := VersionSourceGit{
			Repository: t.Repository,
			Reference:  tag.Reference,
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Git: &source,
			},
		})
	}
	return out, nil
}

// SourceDocker is a Docker repository used as a version source
type SourceDocker struct {
	Repository *simpledocker.Repository
}

// Fetch retrieves Docker tags as versions
func (t SourceDocker) Fetch() (WithSources, error) {
	tags, err := t.Repository.Tags()
	if err != nil {
		return nil, err
	}
	out := make(WithSources, 0, len(tags))
	for _, tag := range tags {
		version, err := semver.Parse(tag.Name)
		if err != nil {
			continue
		}
		source := VersionSourceDocker{
			Tag:   tag.Name,
			Image: tag.Image,
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Docker: &source,
			},
		})
	}
	return out, nil
}
