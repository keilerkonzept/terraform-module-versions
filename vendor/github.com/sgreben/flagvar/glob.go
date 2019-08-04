package flagvar

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

// Glob is a `flag.Value` for glob syntax arguments.
// By default, `filepath.Separator` is used as a separator.
// If `Separators` is non-nil, its elements are used as separators.
// To have no separators, set `Separators` to a (non-nil) pointer to an empty slice.
type Glob struct {
	Separators *[]rune

	Value glob.Glob
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Glob) Help() string {
	separators := []rune{filepath.Separator}
	if fv.Separators != nil {
		separators = *fv.Separators
	}
	if len(separators) == 0 {
		return "a glob expression"
	}
	if len(separators) == 1 {
		return fmt.Sprintf("a glob expression with separator %q", separators[0])
	}
	return fmt.Sprintf("a glob expression with separators %q", separators)
}

// Set is flag.Value.Set
func (fv *Glob) Set(v string) error {
	separators := fv.Separators
	if separators == nil {
		separators = &[]rune{filepath.Separator}
	}
	g, err := glob.Compile(v, *separators...)
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = g
	return nil
}

func (fv *Glob) String() string {
	return fv.Text
}

// Globs is a `flag.Value` for glob syntax arguments.
// By default, `filepath.Separator` is used as a separator.
// If `Separators` is non-nil, its elements are used as separators.
// To have no separators, set `Separators` to a (non-nil) pointer to an empty slice.
type Globs struct {
	Separators *[]rune

	Values []glob.Glob
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Globs) Help() string {
	separators := []rune{filepath.Separator}
	if fv.Separators != nil {
		separators = *fv.Separators
	}
	if len(separators) == 0 {
		return "a glob expression"
	}
	if len(separators) == 1 {
		return fmt.Sprintf("a glob expression with separator %q", separators[0])
	}
	return fmt.Sprintf("a glob expression with separators %q", separators)
}

// Set is flag.Value.Set
func (fv *Globs) Set(v string) error {
	separators := fv.Separators
	if separators == nil {
		separators = &[]rune{filepath.Separator}
	}
	g, err := glob.Compile(v, *separators...)
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, g)
	return nil
}

func (fv *Globs) String() string {
	return strings.Join(fv.Texts, ",")
}
