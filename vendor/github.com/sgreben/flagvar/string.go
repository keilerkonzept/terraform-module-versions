package flagvar

import (
	"fmt"
	"sort"
	"strings"
)

// Strings is a `flag.Value` for `string` arguments.
type Strings struct {
	Values []string
}

// Set is flag.Value.Set
func (fv *Strings) Set(v string) error {
	fv.Values = append(fv.Values, v)
	return nil
}

func (fv *Strings) String() string {
	return strings.Join(fv.Values, ",")
}

// StringSet is a `flag.Value` for `string` arguments.
// Only distinct values are returned.
type StringSet struct {
	Value map[string]bool
}

// Values returns a string slice of specified values.
func (fv *StringSet) Values() (out []string) {
	for v := range fv.Value {
		out = append(out, v)
	}
	sort.Strings(out)
	return
}

// Set is flag.Value.Set
func (fv *StringSet) Set(v string) error {
	if fv.Value == nil {
		fv.Value = make(map[string]bool)
	}
	fv.Value[v] = true
	return nil
}

func (fv *StringSet) String() string {
	return strings.Join(fv.Values(), ",")
}

// StringSetCSV is a `flag.Value` for comma-separated string arguments.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type StringSetCSV struct {
	Separator  string
	Accumulate bool

	Value  map[string]bool
	Values []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *StringSetCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of strings", separator)
}

// Set is flag.Value.Set
func (fv *StringSetCSV) Set(v string) error {
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if !fv.Accumulate || fv.Value == nil {
		fv.Value = make(map[string]bool)
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if fv.Value[part] {
			continue
		}
		fv.Value[part] = true
		fv.Values = append(fv.Values, part)
	}
	return nil
}

func (fv *StringSetCSV) String() string {
	return strings.Join(fv.Values, ",")
}
