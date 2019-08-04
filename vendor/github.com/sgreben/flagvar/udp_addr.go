package flagvar

import (
	"fmt"
	"strings"

	"net"
)

// UDPAddr is a `flag.Value` for UDP addresses.
// The `Network` field is used if set, otherwise "udp".
type UDPAddr struct {
	Network string

	Value *net.UDPAddr
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UDPAddr) Help() string {
	return "a UDP address"
}

// Set is flag.Value.Set
func (fv *UDPAddr) Set(v string) error {
	network := "udp"
	if fv.Network != "" {
		network = fv.Network
	}
	udpAddr, err := net.ResolveUDPAddr(network, v)
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = udpAddr
	return nil
}

func (fv *UDPAddr) String() string {
	return fv.Text
}

// UDPAddrs is a `flag.Value` for UDPAddr addresses.
// The `Network` field is used if set, otherwise "udp".
type UDPAddrs struct {
	Network string

	Values []*net.UDPAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UDPAddrs) Help() string {
	return "a UDP address"
}

// Set is flag.Value.Set
func (fv *UDPAddrs) Set(v string) error {
	network := "udp"
	if fv.Network != "" {
		network = fv.Network
	}
	udpAddr, err := net.ResolveUDPAddr(network, v)
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, udpAddr)
	return nil
}

func (fv *UDPAddrs) String() string {
	return strings.Join(fv.Texts, ",")
}

// UDPAddrsCSV is a `flag.Value` for UDPAddr addresses.
// The `Network` field is used if set, otherwise "udp".
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
type UDPAddrsCSV struct {
	Network    string
	Separator  string
	Accumulate bool

	Values []*net.UDPAddr
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *UDPAddrsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of UDP addresses", separator)
}

// Set is flag.Value.Set
func (fv *UDPAddrsCSV) Set(v string) error {
	network := "udp"
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
		udpAddr, err := net.ResolveUDPAddr(network, part)
		if err != nil {
			return err
		}
		fv.Texts = append(fv.Texts, part)
		fv.Values = append(fv.Values, udpAddr)
	}
	return nil
}

func (fv *UDPAddrsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
