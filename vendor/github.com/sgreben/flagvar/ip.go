package flagvar

import (
	"fmt"
	"strings"

	"net"
)

// IP is a `flag.Value` for IP addresses.
type IP struct {
	Value net.IP
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *IP) Help() string {
	return "an IP address"
}

// Set is flag.Value.Set
func (fv *IP) Set(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return fmt.Errorf(`not a valid IP address: "%s"`, v)
	}
	fv.Text = v
	fv.Value = ip
	return nil
}

func (fv *IP) String() string {
	return fv.Text
}

// IPs is a `flag.Value` for IP addresses.
type IPs struct {
	Values []net.IP
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *IPs) Help() string {
	return "an IP address"
}

// Set is flag.Value.Set
func (fv *IPs) Set(v string) error {
	ip := net.ParseIP(v)
	if ip == nil {
		return fmt.Errorf(`not a valid IP address: "%s"`, v)
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, ip)
	return nil
}

func (fv *IPs) String() string {
	return strings.Join(fv.Texts, ",")
}

// IPsCSV is a `flag.Value` for IP addresses.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
type IPsCSV struct {
	Separator  string
	Accumulate bool

	Values []net.IP
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *IPsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of IP addresses", separator)
}

// Set is flag.Value.Set
func (fv *IPsCSV) Set(v string) error {
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
		ip := net.ParseIP(part)
		if ip == nil {
			return fmt.Errorf(`not a valid IP address: "%s"`, part)
		}
		fv.Texts = append(fv.Texts, part)
		fv.Values = append(fv.Values, ip)
	}
	return nil
}

func (fv *IPsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
