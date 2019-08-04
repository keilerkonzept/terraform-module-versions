package flagvar

import (
	"fmt"
	"strconv"
	"strings"
)

// Float is a `flag.Value` for a float argument.
// The `BitSize` field is used for parsing when set.
type Float struct {
	BitSize int

	Value float64
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Float) Help() string {
	var bitSize string
	if fv.BitSize != 0 {
		bitSize = fmt.Sprintf("%d-bit ", fv.BitSize)
	}
	return fmt.Sprintf("a %sfloat", bitSize)
}

// Set is flag.Value.Set
func (fv *Float) Set(v string) error {
	bitSize := fv.BitSize
	if bitSize == 0 {
		bitSize = 64
	}
	n, err := strconv.ParseFloat(v, bitSize)
	if err == nil {
		fv.Value = n
		fv.Text = v
	}
	return err
}

func (fv *Float) String() string {
	return fv.Text
}

// Floats is a `flag.Value` for float arguments.
// The `BitSize` field is used for parsing when set.
type Floats struct {
	BitSize int

	Values []float64
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Floats) Help() string {
	var bitSize string
	if fv.BitSize != 0 {
		bitSize = fmt.Sprintf("%d-bit ", fv.BitSize)
	}
	return fmt.Sprintf("a %sfloat", bitSize)
}

// Set is flag.Value.Set
func (fv *Floats) Set(v string) error {
	bitSize := fv.BitSize
	if bitSize == 0 {
		bitSize = 64
	}
	n, err := strconv.ParseFloat(v, bitSize)
	if err == nil {
		fv.Values = append(fv.Values, n)
		fv.Texts = append(fv.Texts, v)
	}
	return err
}

func (fv *Floats) String() string {
	return strings.Join(fv.Texts, ",")
}

// FloatsCSV is a `flag.Value` for comma-separated `float` arguments.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `BitSize` fields are used for parsing when set.
// The `Separator` field is used instead of the comma when set.
type FloatsCSV struct {
	BitSize    int
	Separator  string
	Accumulate bool

	Values []float64
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *FloatsCSV) Help() string {
	var bitSize string
	if fv.BitSize != 0 {
		bitSize = fmt.Sprintf("%d-bit ", fv.BitSize)
	}
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of %sfloats", separator, bitSize)
}

// Set is flag.Value.Set
func (fv *FloatsCSV) Set(v string) error {
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
		n, err := strconv.ParseFloat(part, bitSize)
		if err != nil {
			return err
		}
		fv.Values = append(fv.Values, n)
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *FloatsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
