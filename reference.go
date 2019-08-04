package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

const (
	terraformModuleSourceGitPrefix         = "git::"
	terraformModuleSourceGithubHTTPSPrefix = "github.com/"
	terraformModuleSourceGithubSSHPrefix   = "git@github.com:"
	terraformRegistryHostname              = "registry.terraform.io"
)

type moduleSource struct {
	Source   string
	Version  *string
	Git      *moduleReferenceGit
	Registry *moduleReferenceRegistry
}

func (s *moduleSource) InferredVersion() *string {
	if s.Git != nil {
		if s.Git.RefValue != nil {
			return s.Git.RefValue
		}
		return nil
	}
	return nil
}

func (s *moduleSource) Type() string {
	if s.Git != nil {
		return "git"
	}
	if s.Registry != nil {
		return "terraform-registry"
	}
	return "unknown"
}

type moduleReference struct {
	Name     string
	Path     string
	Source   string
	Version  *string
	Git      *moduleReferenceGit
	Registry *moduleReferenceRegistry
}

func (r *moduleReference) SourceStruct() moduleSource {
	return moduleSource{r.Source, r.Version, r.Git, r.Registry}
}

type moduleReferenceRegistry struct {
	Hostname  string
	Namespace string
	Name      string
	Provider  string
}

type moduleReferenceGit struct {
	Remote     string
	RefValue   *string
	RemotePath *string
}

func (r *moduleReference) EncodeGit() {
	if r.Git == nil {
		return
	}
	source := terraformModuleSourceGitPrefix + r.Git.Remote
	if r.Git.RemotePath != nil {
		source = source + "//" + *r.Git.RemotePath
	}
	if r.Git.RefValue != nil {
		source = source + "?ref=" + *r.Git.RefValue
	}
	r.Source = source
	r.Version = nil
	r.Git = nil
}

func (r *moduleReference) ParseSource() error {
	sourcePartsNum := strings.Count(r.Source, "/")
	switch {
	case strings.HasPrefix(r.Source, terraformModuleSourceGitPrefix),
		strings.HasPrefix(r.Source, terraformModuleSourceGithubHTTPSPrefix),
		strings.HasPrefix(r.Source, terraformModuleSourceGithubSSHPrefix):
		return r.parseGit()
	case sourcePartsNum == 2:
		return r.parsePublicRegistry()
	case sourcePartsNum == 3:
		return r.parsePrivateRegistry()
	}
	return nil
}

func (r *moduleReference) parsePublicRegistry() error {
	sourceParts := strings.Split(r.Source, "/")
	if len(sourceParts) != 3 {
		return fmt.Errorf("not a public registry module: %q", r.Source)
	}
	r.Registry = &moduleReferenceRegistry{
		Hostname:  terraformRegistryHostname,
		Namespace: sourceParts[0],
		Name:      sourceParts[1],
		Provider:  sourceParts[2],
	}
	return nil
}

func (r *moduleReference) parsePrivateRegistry() error {
	sourceParts := strings.Split(r.Source, "/")
	if len(sourceParts) != 4 {
		return fmt.Errorf("not a private registry module: %q", r.Source)
	}
	r.Registry = &moduleReferenceRegistry{
		Hostname:  sourceParts[0],
		Namespace: sourceParts[1],
		Name:      sourceParts[2],
		Provider:  sourceParts[3],
	}
	return nil
}

func (r *moduleReference) parseGit() error {
	sourceURL, err := url.Parse(strings.TrimPrefix(r.Source, terraformModuleSourceGitPrefix))
	if err != nil {
		return err
	}
	var gitMeta moduleReferenceGit
	if refValue := sourceURL.Query().Get("ref"); refValue != "" {
		gitMeta.RefValue = &refValue
		query := sourceURL.Query()
		query.Del("ref")
		sourceURL.RawQuery = query.Encode()
	}
	gitMeta.Remote = sourceURL.String()
	sourceURLPathParts := strings.Split(sourceURL.Path, "/")
	for i := len(sourceURLPathParts) - 1; i >= 1; i-- {
		if sourceURLPathParts[i] == "" {
			remoteURL := *sourceURL
			remoteURL.Path = path.Join(sourceURLPathParts[:i]...)
			gitMeta.Remote = remoteURL.String()
			remotePath := path.Join(sourceURLPathParts[i+1:]...)
			gitMeta.RemotePath = &remotePath
			break
		}
	}
	r.Git = &gitMeta
	return nil
}

func (r moduleReference) String() string {
	return fmt.Sprintf("%q (%q in %s)", r.Source, r.Name, r.Path)
}
