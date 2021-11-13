package update

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/keilerkonzept/terraform-module-versions/pkg/registry"
	"github.com/keilerkonzept/terraform-module-versions/pkg/source"
	"github.com/keilerkonzept/terraform-module-versions/pkg/versions"
)

type Client struct {
	Registry      registry.Client
	GitAuth       transport.AuthMethod
	VersionsCache map[string][]*semver.Version
}

type Update struct {
	LatestMatchingVersion string
	LatestOverallVersion  string
	LatestMatchingUpdate  string
	LatestOverallUpdate   string
}

func (c *Client) Update(s source.Source, current *semver.Version, constraints *semver.Constraints, includePrerelease bool) (*Update, error) {
	versions, err := c.Versions(s)
	if err != nil {
		return nil, err
	}
	var out Update
	for _, v := range versions {
		if !includePrerelease && v.Prerelease() != "" {
			continue
		}
		versionString := v.Original()
		out.LatestOverallVersion = versionString
		if current != nil && !v.GreaterThan(current) {
			continue
		}
		out.LatestOverallUpdate = versionString
		if constraints == nil || !constraints.Check(v) {
			continue
		}
		out.LatestMatchingVersion = versionString
		if current != nil {
			out.LatestMatchingUpdate = versionString
		}
	}
	return &out, nil
}

func (c *Client) Versions(s source.Source) ([]*semver.Version, error) {
	if c.VersionsCache == nil {
		c.VersionsCache = make(map[string][]*semver.Version, 1)
	}
	if versions, ok := c.VersionsCache[s.URI()]; ok {
		return versions, nil
	}
	switch {
	case s.Git != nil:
		git := s.Git
		versions, err := versions.Git(git.Remote, c.GitAuth)
		if err != nil {
			return nil, fmt.Errorf("fetch versions from %q: %w", git.Remote, err)
		}
		c.VersionsCache[s.URI()] = versions
		return versions, nil
	case s.Registry != nil:
		reg := s.Registry
		versions, err := versions.Registry(c.Registry, reg.Hostname, reg.Namespace, reg.Name, reg.Provider)
		if err != nil {
			return nil, fmt.Errorf("fetch versions from registry: %w", err)
		}
		c.VersionsCache[s.URI()] = versions
		return versions, nil
	case s.Local != nil:
		return nil, nil
	default:
		return nil, source.ErrSourceNotSupported
	}
}
