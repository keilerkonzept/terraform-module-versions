package flagvar

import (
	"fmt"
	"strings"

	"net"
)

// CIDR is a `flag.Value` for CIDR notation IP address and prefix length arguments.
type CIDR struct {
	Value struct {
		IPNet *net.IPNet
		IP    net.IP
	}
	Text string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *CIDR) Help() string {
	return "a CIDR notation IP address and prefix length"
}

// Set is flag.Value.Set
func (fv *CIDR) Set(v string) error {
	ip, ipNet, err := net.ParseCIDR(v)
	if err != nil {
		return err
	}
	fv.Text = v
	fv.Value = struct {
		IPNet *net.IPNet
		IP    net.IP
	}{IP: ip, IPNet: ipNet}
	return nil
}

func (fv *CIDR) String() string {
	return fv.Text
}

// CIDRs is a `flag.Value` for CIDR notation IP address and prefix length arguments.
type CIDRs struct {
	Values []struct {
		IPNet *net.IPNet
		IP    net.IP
	}
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *CIDRs) Help() string {
	return "a CIDR notation IP address and prefix length"
}

// Set is flag.Value.Set
func (fv *CIDRs) Set(v string) error {
	ip, ipNet, err := net.ParseCIDR(v)
	if err != nil {
		return err
	}
	fv.Texts = append(fv.Texts, v)
	fv.Values = append(fv.Values, struct {
		IPNet *net.IPNet
		IP    net.IP
	}{IP: ip, IPNet: ipNet})
	return nil
}

func (fv *CIDRs) String() string {
	return strings.Join(fv.Texts, ",")
}

// CIDRsCSV is a `flag.Value` for CIDR notation IP address and prefix length arguments.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
type CIDRsCSV struct {
	Separator  string
	Accumulate bool

	Values []struct {
		IPNet *net.IPNet
		IP    net.IP
	}
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *CIDRsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	return fmt.Sprintf("%q-separated list of CIDR notation IP addresses/prefix lengths", separator)
}

// Set is flag.Value.Set
func (fv *CIDRsCSV) Set(v string) error {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	if !fv.Accumulate {
		fv.Values = fv.Values[:0]
		fv.Texts = fv.Texts[:0]
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		ip, ipNet, err := net.ParseCIDR(part)
		if err != nil {
			return err
		}
		fv.Texts = append(fv.Texts, part)
		fv.Values = append(fv.Values, struct {
			IPNet *net.IPNet
			IP    net.IP
		}{IP: ip, IPNet: ipNet})
	}
	return nil
}

func (fv *CIDRsCSV) String() string {
	return strings.Join(fv.Texts, ",")
}
