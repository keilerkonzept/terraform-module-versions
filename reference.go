package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Module source prefixes
const (
	ModuleSourceGitSSHPrefix        = "git::ssh://"
	ModuleSourceGitPrefix           = "git::"
	ModuleSourceGithubHTTPSPrefix   = "github.com/"
	ModuleSourceGithubSSHPrefix     = "git@github.com:"
	ModuleSourceLocalPathPrefix1    = "./"
	ModuleSourceLocalPathPrefix2    = "../"
	terraformPublicRegistryHostname = "registry.terraform.io"
)

type moduleSource struct {
	Source    string
	Version   *string
	Git       *moduleReferenceGit
	Registry  *moduleReferenceRegistry
	LocalPath *string
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
	switch {
	case s.LocalPath != nil:
		return "local"
	case s.Git != nil:
		return "git"
	case s.Registry != nil:
		return "registry"
	default:
		return "unknown"
	}
}

type moduleReference struct {
	Name      string
	Path      string
	Source    string
	Version   *string
	Git       *moduleReferenceGit
	Registry  *moduleReferenceRegistry
	LocalPath *string
}

func (r *moduleReference) SourceStruct() moduleSource {
	return moduleSource{r.Source, r.Version, r.Git, r.Registry, r.LocalPath}
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
	source := ModuleSourceGitPrefix + r.Git.Remote
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
	case strings.HasPrefix(r.Source, ModuleSourceLocalPathPrefix1),
		strings.HasPrefix(r.Source, ModuleSourceLocalPathPrefix2):
		r.parseLocalPath(r.Source)
		return nil
	case strings.HasPrefix(r.Source, ModuleSourceGitSSHPrefix):
		return r.parseGit(r.Source)
	case strings.HasPrefix(r.Source, ModuleSourceGitPrefix):
	    source := strings.TrimPrefix(r.Source, ModuleSourceGitPrefix)
	    source = strings.Replace(source, ":", "/", 1)
	    return r.parseGit("ssh://" + source)
	case strings.HasPrefix(r.Source, ModuleSourceGithubHTTPSPrefix):
		return r.parseGit("https://" + r.Source)
	case strings.HasPrefix(r.Source, ModuleSourceGithubSSHPrefix):
		source := strings.TrimPrefix(r.Source, ModuleSourceGithubSSHPrefix)
		source = fmt.Sprintf("git@github.com/%s", source)
		return r.parseGit("ssh://" + source)
	case sourcePartsNum == 2:
		return r.parsePublicRegistry(r.Source)
	case sourcePartsNum > 2:
		return r.parsePrivateRegistry(r.Source)
	}
	return nil
}

func (r *moduleReference) parseLocalPath(source string) {
	r.LocalPath = &source
}

func (r *moduleReference) parsePublicRegistry(source string) error {
	sourceParts := strings.Split(source, "/")
	if len(sourceParts) < 3 {
		return fmt.Errorf("not a public registry module: %q", source)
	}
	r.Registry = &moduleReferenceRegistry{
		Hostname:  terraformPublicRegistryHostname,
		Namespace: sourceParts[0],
		Name:      sourceParts[1],
		Provider:  sourceParts[2],
	}
	return nil
}

func (r *moduleReference) parsePrivateRegistry(source string) error {
	sourceParts := strings.Split(source, "/")
	if len(sourceParts) < 4 {
		return fmt.Errorf("not a private registry module: %q", source)
	}
	r.Registry = &moduleReferenceRegistry{
		Hostname:  sourceParts[0],
		Namespace: sourceParts[1],
		Name:      sourceParts[2],
		Provider:  sourceParts[3],
	}
	return nil
}

func (r *moduleReference) parseGit(source string) error {
	sourceURL, err := url.Parse(strings.TrimPrefix(source, ModuleSourceGitPrefix))
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
