package versions

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/keilerkonzept/terraform-module-versions/pkg/registry"
)

func Registry(client registry.Client, hostname, namespace, name, provider string) ([]*semver.Version, error) {
	baseURL, err := client.Discover(hostname)
	if err != nil {
		return nil, fmt.Errorf("discover registry at %q: %w", hostname, err)
	}
	baseURLStruct, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse module registry url %q: %w", baseURL, err)
	}
	if baseURLStruct.Scheme == "" {
		baseURLStruct.Scheme = "https"
	}
	if baseURLStruct.Host == "" {
		baseURLStruct.Host = hostname
	}
	baseURL = baseURLStruct.String()
	versions, err := client.ListVersions(baseURL, namespace, name, provider)
	if err != nil {
		return nil, fmt.Errorf("list versions: %w", err)
	}
	out := make([]*semver.Version, len(versions))
	for i, versionString := range versions {
		version, err := semver.NewVersion(versionString)
		if err != nil {
			version = nil
		}
		out[i] = version
	}
	sort.Sort(semver.Collection(out))
	return out, nil
}
