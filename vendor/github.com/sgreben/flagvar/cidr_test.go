package flagvar_test

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestCIDR(t *testing.T) {
	fv := flagvar.CIDR{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr", "")

	err := fs.Parse([]string{"-cidr", "192.168.0.3/24"})
	if err != nil {
		t.Fail()
	}
	ip, ipNet, _ := net.ParseCIDR("192.168.0.3/24")
	cidr := struct {
		IPNet *net.IPNet
		IP    net.IP
	}{IP: ip, IPNet: ipNet}
	if !reflect.DeepEqual(fv.Value, cidr) {
		t.Fail()
	}
}

func TestCIDRFail(t *testing.T) {
	fv := flagvar.CIDR{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr", "")

	err := fs.Parse([]string{"-cidr", "999.999.999.999/123"})
	if err == nil {
		t.Fail()
	}
}

func TestCIDRs(t *testing.T) {
	fv := flagvar.CIDRs{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr", "")

	err := fs.Parse([]string{"-cidr", "192.168.0.3/24"})
	if err != nil {
		t.Fail()
	}
	ip, ipNet, _ := net.ParseCIDR("192.168.0.3/24")
	if !reflect.DeepEqual(fv.Values, []struct {
		IPNet *net.IPNet
		IP    net.IP
	}{{IP: ip, IPNet: ipNet}}) {
		t.Fail()
	}
}

func TestCIDRsFail(t *testing.T) {
	fv := flagvar.CIDRs{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr", "")

	err := fs.Parse([]string{"-cidr", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestCIDRsCSV(t *testing.T) {
	fv := flagvar.CIDRsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr-csv", "")

	err := fs.Parse([]string{"-cidr-csv", "192.168.0.1/16,10.10.10.10/24,10.10.10.10/24"})
	if err != nil {
		t.Fail()
	}
	ip1, ipNet1, _ := net.ParseCIDR("192.168.0.1/16")
	ip2, ipNet2, _ := net.ParseCIDR("10.10.10.10/24")
	ip3, ipNet3, _ := net.ParseCIDR("10.10.10.10/24")
	if !reflect.DeepEqual(fv.Values, []struct {
		IPNet *net.IPNet
		IP    net.IP
	}{{IP: ip1, IPNet: ipNet1}, {IP: ip2, IPNet: ipNet2}, {IP: ip3, IPNet: ipNet3}}) {
		t.Fail()
	}
}

func TestCIDRsCSVSeparator(t *testing.T) {
	fv := flagvar.CIDRsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr-csv", "")

	err := fs.Parse([]string{"-cidr-csv", "192.168.0.1/16;10.10.10.10/24;10.10.10.10/24"})
	if err != nil {
		t.Fail()
	}
	ip1, ipNet1, _ := net.ParseCIDR("192.168.0.1/16")
	ip2, ipNet2, _ := net.ParseCIDR("10.10.10.10/24")
	ip3, ipNet3, _ := net.ParseCIDR("10.10.10.10/24")
	if !reflect.DeepEqual(fv.Values, []struct {
		IPNet *net.IPNet
		IP    net.IP
	}{{IP: ip1, IPNet: ipNet1}, {IP: ip2, IPNet: ipNet2}, {IP: ip3, IPNet: ipNet3}}) {
		t.Fail()
	}
}

func TestCIDRsCSVAccumulate(t *testing.T) {
	fv := flagvar.CIDRsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr-csv", "")

	err := fs.Parse([]string{"-cidr-csv", "192.168.0.1/16,10.10.10.10/24", "-cidr-csv", "10.10.10.10/24"})
	if err != nil {
		t.Fail()
	}
	ip1, ipNet1, _ := net.ParseCIDR("192.168.0.1/16")
	ip2, ipNet2, _ := net.ParseCIDR("10.10.10.10/24")
	ip3, ipNet3, _ := net.ParseCIDR("10.10.10.10/24")
	if !reflect.DeepEqual(fv.Values, []struct {
		IPNet *net.IPNet
		IP    net.IP
	}{{IP: ip1, IPNet: ipNet1}, {IP: ip2, IPNet: ipNet2}, {IP: ip3, IPNet: ipNet3}}) {
		t.Fail()
	}
}

func TestCIDRsCSVFail(t *testing.T) {
	fv := flagvar.CIDRsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "cidr-csv", "")

	err := fs.Parse([]string{"-cidr-csv", "xxx"})
	if err == nil {
		t.Fail()
	}
}
