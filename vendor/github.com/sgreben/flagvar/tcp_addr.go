package flagvar

import (
	"fmt"
	"strings"

	"net"
)

// TCPAddr is a `flag.Value` for TCP addresses.
// The `Network` field is used if set, otherwise "tcp".
type TCPAddr struct {
	Network string

	Value *net.TCPAddr
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TCPAddr) Help() string {
	return "a TCP address"
}

// Set is flag.Value.Set
func (fv *TCPAddr) Set(v string) error {
	network := "tcp"
	if fv.Network != "" {
		network = fv.Network
	}
	tcpAddr, err := net.ResolveTCPAddr(network, v)
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = tcpAddr
	return nil
}

func (fv *TCPAddr) String() string {
	return fv.Text
}

// TCPAddrs is a `flag.Value` for TCPAddr addresses.
// The `Network` field is used if set, otherwise "tcp".
type TCPAddrs struct {
	Network string

	Values []*net.TCPAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TCPAddrs) Help() string {
	return "a TCP address"
}

// Set is flag.Value.Set
func (fv *TCPAddrs) Set(v string) error {
	network := "tcp"
	if fv.Network != "" {
		network = fv.Network
	}
	tcpAddr, err := net.ResolveTCPAddr(network, v)
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, tcpAddr)
	return nil
}

func (fv *TCPAddrs) String() string {
	return strings.Join(fv.Texts, ",")
}

// TCPAddrsCSV is a `flag.Value` for TCPAddr addresses.
// The `Network` field is used if set, otherwise "tcp".
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
type TCPAddrsCSV struct {
	Network    string
	Separator  string
	Accumulate bool

	Values []*net.TCPAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TCPAddrsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of TCP addresses", separator)
}

// Set is flag.Value.Set
func (fv *TCPAddrsCSV) Set(v string) error {
	network := "tcp"
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
		tcpAddr, err := net.ResolveTCPAddr(network, part)
		if err != nil {
			return err
		}
		fv.Texts = append(fv.Texts, part)
		fv.Values = append(fv.Values, tcpAddr)
	}
	return nil
}

func (fv *TCPAddrsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
