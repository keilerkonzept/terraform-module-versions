package flagvar

import (
	"fmt"
	"strings"
	"time"
)

// Time is a `flag.Value` for `time.Time` arguments.
// The value of the `Layout` field is used for parsing when specified.
// Otherwise, `time.RFC3339` is used.
type Time struct {
	Layout string

	Value time.Time
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Time) Help() string {
	layout := time.RFC3339
	if fv.Layout != "" {
		layout = fv.Layout
	}
	return fmt.Sprintf("a time, e.g. %s", layout)
}

// Set is flag.Value.Set
func (fv *Time) Set(v string) error {
	layout := fv.Layout
	if layout == "" {
		layout = time.RFC3339
	}
	t, err := time.Parse(layout, v)
	if err == nil {
		fv.Text = v
		fv.Value = t
	}
	return err
}

func (fv *Time) String() string {
	return fv.Text
}

// Times is a `flag.Value` for `time.Time` arguments.
// The value of the `Layout` field is used for parsing when specified.
// Otherwise, `time.RFC3339` is used.
type Times struct {
	Layout string

	Values []time.Time
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Times) Help() string {
	layout := time.RFC3339
	if fv.Layout != "" {
		layout = fv.Layout
	}
	return fmt.Sprintf("a time, e.g. %s", layout)
}

// Set is flag.Value.Set
func (fv *Times) Set(v string) error {
	layout := fv.Layout
	if layout == "" {
		layout = time.RFC3339
	}
	t, err := time.Parse(layout, v)
	if err == nil {
		fv.Texts = append(fv.Texts, v)
		fv.Values = append(fv.Values, t)
	}
	return err
}

func (fv *Times) String() string {
	return strings.Join(fv.Texts, ",")
}
