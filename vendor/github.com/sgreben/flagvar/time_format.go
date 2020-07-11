package flagvar

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TimeFormat is a `flag.Value` for `time.Time` formats.
type TimeFormat struct {
	Value string
	Text  string
}

var timeFormats = map[string]string{
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"ISO8601":     "2006-01-02T15:04:05Z07:00",
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TimeFormat) Help() string {
	formats := make([]string, 0, len(timeFormats))
	for k, v := range timeFormats {
		formats = append(formats, fmt.Sprintf("%s (%q)", k, v))
	}
	sort.Strings(formats)
	return fmt.Sprintf("a time format, one of %s", strings.Join(formats, ", "))
}

// Set is flag.Value.Set
func (fv *TimeFormat) Set(v string) error {
	namedFormat, ok := timeFormats[v]
	if !ok {
		return fmt.Errorf("no such time format defined: %q", v)
	}
	fv.Text = v
	fv.Value = namedFormat
	return nil
}

func (fv *TimeFormat) String() string {
	return fv.Text
}
