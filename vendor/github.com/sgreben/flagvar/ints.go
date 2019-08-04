package flagvar

import (
	"fmt"
	"strconv"
	"strings"
)

// Ints is a `flag.Value` for `int` arguments.
// The `Base` and `BitSize` fields are used for parsing when set.
type Ints struct {
	Base    int
	BitSize int

	Values []int64
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Ints) Help() string {
	var base, bitSize string
	if fv.Base != 0 {
		base = fmt.Sprintf("base %d ", fv.Base)
	}
	if fv.BitSize != 0 {
		bitSize = fmt.Sprintf("%d-bit ", fv.BitSize)
	}
	if base != "" || bitSize != "" {
		return fmt.Sprintf("a %s%sinteger", bitSize, base)
	}
	return "an integer"
}

// Set is flag.Value.Set
func (fv *Ints) Set(v string) error {
	base := fv.Base
	if base == 0 {
		base = 10
	}
	bitSize := fv.BitSize
	if bitSize == 0 {
		bitSize = 64
	}
	n, err := strconv.ParseInt(v, base, bitSize)
	if err == nil {
		fv.Values = append(fv.Values, n)
		fv.Texts = append(fv.Texts, v)
	}
	return err
}

func (fv *Ints) String() string {
	return strings.Join(fv.Texts, ",")
}

// IntsCSV is a `flag.Value` for comma-separated `int` arguments.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Base` and `BitSize` fields are used for parsing when set.
// The `Separator` field is used instead of the comma when set.
type IntsCSV struct {
	Base       int
	BitSize    int
	Separator  string
	Accumulate bool

	Values []int64
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *IntsCSV) Help() string {
	var base, bitSize string
	if fv.Base != 0 {
		base = fmt.Sprintf("base %d ", fv.Base)
	}
	if fv.BitSize != 0 {
		bitSize = fmt.Sprintf("%d-bit ", fv.BitSize)
	}
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of %s%sintegers", separator, bitSize, base)
}

// Set is flag.Value.Set
func (fv *IntsCSV) Set(v string) error {
	base := fv.Base
	if base == 0 {
		base = 10
	}
	bitSize := fv.BitSize
	if bitSize == 0 {
		bitSize = 64
	}
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if !fv.Accumulate {
		fv.Values = fv.Values[:0]
		fv.Texts = fv.Texts[:0]
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		n, err := strconv.ParseInt(part, base, bitSize)
		if err != nil {
			return err
		}
		fv.Values = append(fv.Values, n)
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *IntsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
