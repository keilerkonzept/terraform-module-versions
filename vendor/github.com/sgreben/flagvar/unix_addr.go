package flagvar

import (
	"fmt"
	"strings"

	"net"
)

// UnixAddr is a `flag.Value` for Unix addresses.
// The `Network` field is used if set, otherwise "unix".
type UnixAddr struct {
	Network string

	Value *net.UnixAddr
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UnixAddr) Help() string {
	return "a UNIX domain socket address"
}

// Set is flag.Value.Set
func (fv *UnixAddr) Set(v string) error {
	network := "unix"
	if fv.Network != "" {
		network = fv.Network
	}
	unixAddr, err := net.ResolveUnixAddr(network, v)
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = unixAddr
	return nil
}

func (fv *UnixAddr) String() string {
	return fv.Text
}

// UnixAddrs is a `flag.Value` for UnixAddr addresses.
// The `Network` field is used if set, otherwise "unix".
type UnixAddrs struct {
	Network string

	Values []*net.UnixAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UnixAddrs) Help() string {
	return "a UNIX domain socket address"
}

// Set is flag.Value.Set
func (fv *UnixAddrs) Set(v string) error {
	network := "unix"
	if fv.Network != "" {
		network = fv.Network
	}
	unixAddr, err := net.ResolveUnixAddr(network, v)
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, unixAddr)
	return nil
}

func (fv *UnixAddrs) String() string {
	return strings.Join(fv.Texts, ",")
}

// UnixAddrsCSV is a `flag.Value` for UnixAddr addresses.
// The `Network` field is used if set, otherwise "unix".
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
type UnixAddrsCSV struct {
	Network    string
	Separator  string
	Accumulate bool

	Values []*net.UnixAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UnixAddrsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of UNIX domain socket addresses", separator)
}

// Set is flag.Value.Set
func (fv *UnixAddrsCSV) Set(v string) error {
	network := "unix"
	if fv.Network != "" {
		network = fv.Network
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
		unixAddr, err := net.ResolveUnixAddr(network, part)
		if err != nil {
			return err
		}
		fv.Texts = append(fv.Texts, part)
		fv.Values = append(fv.Values, unixAddr)
	}
	return nil
}

func (fv *UnixAddrsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
