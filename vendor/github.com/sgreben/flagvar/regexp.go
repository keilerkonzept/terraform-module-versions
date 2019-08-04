package flagvar

import (
	"regexp"
	"strings"
)

// Regexp is a `flag.Value` for regular expression arguments.
// If `POSIX` is set to true, `regexp.CompilePOSIX` is used instead of `regexp.Compile`.
type Regexp struct {
	POSIX bool

	Value *regexp.Regexp
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Regexp) Help() string {
	if fv.POSIX {
		return "a POSIX regular expression"
	}
	return "a regular expression"
}

// Set is flag.Value.Set
func (fv *Regexp) Set(v string) error {
	var err error
	var re *regexp.Regexp
	if fv.POSIX {
		re, err = regexp.CompilePOSIX(v)
	} else {
		re, err = regexp.Compile(v)
	}
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = re
	return nil
}

func (fv *Regexp) String() string {
	return fv.Text
}

// Regexps is a `flag.Value` for regular expression arguments.
type Regexps struct {
	POSIX bool

	Values []*regexp.Regexp
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Regexps) Help() string {
	if fv.POSIX {
		return "a POSIX regular expression"
	}
	return "a regular expression"
}

// Set is flag.Value.Set
func (fv *Regexps) Set(v string) error {
	var err error
	var re *regexp.Regexp
	if fv.POSIX {
		re, err = regexp.CompilePOSIX(v)
	} else {
		re, err = regexp.Compile(v)
	}
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, re)
	return nil
}

func (fv *Regexps) String() string {
	return strings.Join(fv.Texts, ",")
}
