package flagvar_test

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestIP(t *testing.T) {
	fv := flagvar.IP{}
	var fs flag.FlagSet
	fs.Var(&fv, "ip", "")

	err := fs.Parse([]string{"-ip", "127.0.0.1"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, net.ParseIP("127.0.0.1")) {
		t.Fail()
	}
}

func TestIPFail(t *testing.T) {
	fv := flagvar.IP{}
	var fs flag.FlagSet
	fs.Var(&fv, "ip", "")

	err := fs.Parse([]string{"-ip", "999.999.999.999"})
	if err == nil {
		t.Fail()
	}
}

func TestIPs(t *testing.T) {
	fv := flagvar.IPs{}
	var fs flag.FlagSet
	fs.Var(&fv, "ip", "")

	err := fs.Parse([]string{"-ip", "127.0.0.1"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []net.IP{net.ParseIP("127.0.0.1")}) {
		t.Fail()
	}
}

func TestIPsFail(t *testing.T) {
	fv := flagvar.IPs{}
	var fs flag.FlagSet
	fs.Var(&fv, "ip", "")

	err := fs.Parse([]string{"-ip", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestIPsCSV(t *testing.T) {
	fv := flagvar.IPsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "ips-csv", "")

	err := fs.Parse([]string{"-ips-csv", "127.0.0.1,10.10.1.2"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("10.10.1.2")}) {
		t.Fail()
	}
}

func TestIPsCSVSeparator(t *testing.T) {
	fv := flagvar.IPsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "ips-csv", "")

	err := fs.Parse([]string{"-ips-csv", "127.0.0.1;10.10.1.2"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("10.10.1.2")}) {
		t.Fail()
	}
}

func TestIPsCSVAccumulate(t *testing.T) {
	fv := flagvar.IPsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "ips-csv", "")

	err := fs.Parse([]string{"-ips-csv", "127.0.0.1,10.10.1.2", "-ips-csv", "10.10.1.2"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("10.10.1.2"), net.ParseIP("10.10.1.2")}) {
		t.Fail()
	}
}

func TestIPsCSVFail(t *testing.T) {
	fv := flagvar.IPsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "ips-csv", "")

	err := fs.Parse([]string{"-ips-csv", "999.999.999.1"})
	if err == nil {
		t.Fail()
	}
}
