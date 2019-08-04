package flagvar

import (
	"flag"
	"fmt"
	"strings"
)

// WrapPointer wraps a pointer to a `flag.Value`
// This can be used to switch between different argument parsers.
type WrapPointer struct {
	Value *flag.Value
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *WrapPointer) Help() string {
	if fv.Value == nil {
		return ""
	}
	if helper, ok := (*fv.Value).(interface {
		Help() string
	}); ok {
		return helper.Help()
	}
	return ""
}

// Set is flag.Value.Set
func (fv *WrapPointer) Set(v string) error {
	return (*fv.Value).Set(v)
}

func (fv WrapPointer) String() string {
	if fv.Value == nil || *fv.Value == nil {
		return ""
	}
	return (*fv.Value).String()
}

// WrapFunc wraps a nullary function returning a `flag.Value`
// This can be used to switch between different argument parsers.
type WrapFunc func() flag.Value

// Help returns a string suitable for inclusion in a flag help message.
func (fv WrapFunc) Help() string {
	if fv == nil {
		return ""
	}
	if helper, ok := fv().(interface {
		Help() string
	}); ok {
		return helper.Help()
	}
	return ""
}

// Set is flag.Value.Set
func (fv WrapFunc) Set(v string) error {
	return fv().Set(v)
}

func (fv WrapFunc) String() string {
	if fv == nil {
		return ""
	}
	return fv().String()
}

// Wrap wraps a `flag.Value` and calls `Updated` each time the underlying value is set.
type Wrap struct {
	Value   flag.Value
	Updated func()
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Wrap) Help() string {
	if fv.Value == nil {
		return ""
	}
	if helper, ok := fv.Value.(interface {
		Help() string
	}); ok {
		return helper.Help()
	}
	return ""
}

// Set is flag.Value.Set
func (fv *Wrap) Set(v string) error {
	err := fv.Value.Set(v)
	if err == nil {
		fv.Updated()
	}
	return err
}

func (fv *Wrap) String() string {
	if fv.Value == nil {
		return ""
	}
	return fv.Value.String()
}

// WrapCSV wraps a `flag.Value` and calls `UpdatedOne` after each single value and `UpdatedAll` after each CSV batch.
// The `Separator` field is used instead of the comma when set.
type WrapCSV struct {
	Value      flag.Value
	Separator  string
	UpdatedOne func()
	UpdatedAll func()
	StringFunc func() string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *WrapCSV) Help() string {
	if fv.Value == nil {
		return ""
	}
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if helper, ok := fv.Value.(interface {
		Help() string
	}); ok {
		return fmt.Sprintf("%q-separated values, each value %s", separator, helper.Help())
	}
	return ""
}

// Set is flag.Value.Set
func (fv *WrapCSV) Set(v string) error {
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		err := fv.Value.Set(part)
		if err != nil {
			return err
		}
		if fv.UpdatedOne != nil {
			fv.UpdatedOne()
		}
	}
	if fv.UpdatedAll != nil {
		fv.UpdatedAll()
	}
	return nil
}

func (fv *WrapCSV) String() string {
	if fv.StringFunc != nil {
		return fv.StringFunc()
	}
	if fv.Value == nil {
		return ""
	}
	return fv.Value.String()
}
