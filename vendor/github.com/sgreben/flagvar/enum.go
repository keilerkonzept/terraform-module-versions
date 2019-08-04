package flagvar

import (
	"fmt"
	"sort"
	"strings"
)

// Enum is a `flag.Value` for one-of-a-fixed-set string arguments.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type Enum struct {
	Choices       []string
	CaseSensitive bool

	Value string
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Enum) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Set is flag.Value.Set
func (fv *Enum) Set(v string) error {
	fv.Text = v
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	for _, c := range fv.Choices {
		if equal(c, v) {
			fv.Value = c
			return nil
		}
	}
	return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
}

func (fv *Enum) String() string {
	return fv.Value
}

// Enums is a `flag.Value` for one-of-a-fixed-set string arguments.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type Enums struct {
	Choices       []string
	CaseSensitive bool

	Values []string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Enums) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Set is flag.Value.Set
func (fv *Enums) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	for _, c := range fv.Choices {
		if equal(c, v) {
			fv.Values = append(fv.Values, c)
			fv.Texts = append(fv.Texts, v)
			return nil
		}
	}
	return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
}

func (fv *Enums) String() string {
	return strings.Join(fv.Values, ",")
}

// EnumsCSV is a `flag.Value` for comma-separated enum arguments.
// The value of the `Choices` field defines the valid choices.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumsCSV struct {
	Choices       []string
	Separator     string
	Accumulate    bool
	CaseSensitive bool

	Values []string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	if fv.CaseSensitive {
		return fmt.Sprintf("%q-separated list of values from %v (case-sensitive)", separator, fv.Choices)
	}
	return fmt.Sprintf("%q-separated list of values from %v", separator, fv.Choices)
}

// Set is flag.Value.Set
func (fv *EnumsCSV) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if !fv.Accumulate {
		fv.Values = fv.Values[:0]
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var ok bool
		var value string
		for _, c := range fv.Choices {
			if equal(c, part) {
				value = c
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
		}
		fv.Values = append(fv.Values, value)
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *EnumsCSV) String() string {
	return strings.Join(fv.Values, ",")
}

// EnumSet is a `flag.Value` for one-of-a-fixed-set string arguments.
// Only distinct values are returned.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumSet struct {
	Choices       []string
	CaseSensitive bool

	Value map[string]bool
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumSet) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Values returns a string slice of specified values.
func (fv *EnumSet) Values() (out []string) {
	for v := range fv.Value {
		out = append(out, v)
	}
	sort.Strings(out)
	return
}

// Set is flag.Value.Set
func (fv *EnumSet) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	var ok bool
	for _, c := range fv.Choices {
		if equal(c, v) {
			v = c
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
	}
	if fv.Value == nil {
		fv.Value = make(map[string]bool)
	}
	fv.Value[v] = true
	fv.Texts = append(fv.Texts, v)
	return nil
}

func (fv *EnumSet) String() string {
	return strings.Join(fv.Values(), ",")
}

// EnumSetCSV is a `flag.Value` for comma-separated enum arguments.
// Only distinct values are returned.
// The value of the `Choices` field defines the valid choices.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumSetCSV struct {
	Choices       []string
	Separator     string
	Accumulate    bool
	CaseSensitive bool

	Value map[string]bool
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumSetCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	if fv.CaseSensitive {
		return fmt.Sprintf("%q-separated list of values from %v (case-sensitive)", separator, fv.Choices)
	}
	return fmt.Sprintf("%q-separated list of values from %v", separator, fv.Choices)
}

// Values returns a string slice of specified values.
func (fv *EnumSetCSV) Values() (out []string) {
	for v := range fv.Value {
		out = append(out, v)
	}
	sort.Strings(out)
	return
}

// Set is flag.Value.Set
func (fv *EnumSetCSV) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
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
		var ok bool
		var value string
		for _, c := range fv.Choices {
			if equal(c, part) {
				value = c
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
		}
		fv.Value[value] = true
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *EnumSetCSV) String() string {
	return strings.Join(fv.Values(), ",")
}
