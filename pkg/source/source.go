package source

import (
	"errors"
	"fmt"

	tfaddr "github.com/hashicorp/terraform-registry-address"
)

type Source struct {
	Git      *Git
	Registry *Registry
	Local    *string
}

func (s Source) Type() string {
	switch {
	case s.Git != nil:
		return "git"
	case s.Registry != nil:
		return "registry"
	case s.Local != nil:
		return "local"
	}
	return ""
}

func (s Source) URI() string {
	switch {
	case s.Git != nil:
		return s.Git.Remote
	case s.Registry != nil:
		return s.Registry.Normalized
	case s.Local != nil:
		return *s.Local
	}
	return ""
}

var ErrSourceNotSupported = errors.New("source not supported")

func Parse(raw string) (*Source, error) {
	if module, err := tfaddr.ParseModuleSource(raw); err == nil {
		out := &Source{
			Registry: &Registry{
				Hostname:     module.Package.Host.String(),
				Namespace:    module.Package.Namespace,
				Name:         module.Package.Name,
				TargetSystem: module.Package.TargetSystem,
				Normalized:   module.ForDisplay(),
			},
		}
		return out, nil
	}
	proto, detected, err := detect(raw)
	if err != nil {
		return nil, err
	}
	switch proto {
	case "git":
		git, err := parseGitURL(detected)
		if err != nil {
			return nil, err
		}
		return &Source{Git: git}, nil
	case "file":
		return &Source{Local: &detected}, nil
	default:
		return nil, fmt.Errorf("%w: %v (%v)", ErrSourceNotSupported, proto, raw)
	}
}
