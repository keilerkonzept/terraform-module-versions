package versions

import (
	"fmt"
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Git(remoteURL string, auth transport.AuthMethod) ([]*semver.Version, error) {
	raw, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, fmt.Errorf("git init: %w", err)
	}
	remote, err := raw.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{remoteURL},
	})
	if err != nil {
		return nil, fmt.Errorf("git remote: %w", err)
	}
	if githubToken := os.Getenv("GITHUB_TOKEN"); githubToken != "" {
		auth = &http.BasicAuth{
			Username: githubToken,
		}
	}
	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return nil, fmt.Errorf("git list refs: %w", err)
	}
	out := make([]*semver.Version, 0, len(refs))
	for _, ref := range refs {
		version, err := semver.NewVersion(ref.Name().Short())
		if err != nil {
			continue
		}
		out = append(out, version)
	}
	sort.Sort(semver.Collection(out))
	return out, nil
}
