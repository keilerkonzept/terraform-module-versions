package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"gopkg.in/src-d/go-git.v4"

	"github.com/sgreben/versions/pkg/semver"

	"github.com/sgreben/versions/pkg/simplegit"
	"github.com/sgreben/versions/pkg/versions"
)

type moduleVersion struct {
	VersionString         string
	Version               *semver.Version
	VersionSourceGit      *versions.VersionSourceGit
	VersionSourceRegistry *versionSourceRegistry
}

func (r *moduleSource) Versions() ([]moduleVersion, error) {
	if r.Git != nil {
		return r.versionsGit()
	}
	if r.Registry != nil {
		return r.versionsRegistry()
	}
	return nil, nil
}

func (r *moduleSource) versionsGit() ([]moduleVersion, error) {
	repository := simplegit.Repository{URL: r.Git.Remote, CloneOptions: git.CloneOptions{
		SingleBranch: true,
		Depth:        1,
		Tags:         git.NoTags,
		NoCheckout:   true,
	}}
	vs := versions.SourceGit{Repository: repository}
	versions, err := vs.Fetch()
	if err != nil {
		return nil, fmt.Errorf("%q: %v", r.Source, err)
	}
	sort.Sort(versions)
	var out []moduleVersion
	for _, v := range versions {
		out = append(out, moduleVersion{
			VersionString:    v.Version.String(),
			Version:          v.Version,
			VersionSourceGit: v.Source.Git,
		})
	}
	return out, nil
}

func (r *moduleSource) versionsRegistry() ([]moduleVersion, error) {
	var client http.Client
	baseURL, err := registryDiscover(&client, r.Registry.Hostname)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", r.Source, err)
	}
	baseURLStruct, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", r.Source, err)
	}
	if baseURLStruct.Scheme == "" {
		baseURLStruct.Scheme = "https"
	}
	if baseURLStruct.Host == "" {
		baseURLStruct.Host = r.Registry.Hostname
	}
	baseURL = baseURLStruct.String()
	versions, err := registryListVersions(&client, baseURL, r.Registry.Namespace, r.Registry.Name, r.Registry.Provider)
	if err != nil {
		return nil, err
	}
	var moduleVersions []moduleVersion
	for _, versionString := range versions {
		version, err := semver.Parse(versionString)
		if err != nil {
			version = nil
		}
		moduleVersions = append(moduleVersions, moduleVersion{
			VersionString: versionString,
			Version:       version,
			VersionSourceRegistry: &versionSourceRegistry{
				Hostname: r.Registry.Hostname,
			},
		})
	}
	return moduleVersions, nil
}

type versionSourceRegistry struct {
	Hostname string
}
