package simplegit

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Repository represents a git repository
type Repository struct {
	URL          string           `json:"URL"`
	CloneOptions git.CloneOptions `json:"-"`
	Cached       *git.Repository  `json:"-"`
}

// Fetch clones the repository to an in-memory filesystem
// The CloneOptions field is used for cloning.
func (r *Repository) Fetch() error {
	s := memory.NewStorage()
	var fs billy.Filesystem
	if r.CloneOptions.Depth > 0 {
		fs = memfs.New()
	}
	r.CloneOptions.URL = r.URL
	raw, err := git.Clone(s, fs, &r.CloneOptions)
	if err != nil {
		return err
	}
	r.Cached = raw
	return nil
}

// Raw returns a go-git Repository structure for this repository.
// Calls Fetch if necessary.
func (r *Repository) Raw() (raw *git.Repository, err error) {
	if r.Cached == nil {
		err = r.Fetch()
	}
	raw = r.Cached
	return
}

// Tags returns the list of tags in this repository.
// Calls Fetch if necessary.
func (r *Repository) Tags() (out []struct {
	Name      string
	Reference string
}, err error) {
	raw, err := r.Raw()
	if err != nil {
		return nil, err
	}
	iter, err := raw.Tags()
	if err != nil {
		return nil, err
	}
	iter.ForEach(func(tag *plumbing.Reference) (_ error) {
		out = append(out, struct {
			Name      string
			Reference string
		}{
			Name:      tag.Name().Short(),
			Reference: tag.Name().String(),
		})
		return
	})
	return out, nil
}
