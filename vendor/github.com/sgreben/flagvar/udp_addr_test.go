package flagvar_test

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestUDPAddr(t *testing.T) {
	fv := flagvar.UDPAddr{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "127.0.0.1:123"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:123")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestUDPAddrNetwork(t *testing.T) {
	fv := flagvar.UDPAddr{Network: "udp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "127.0.0.1:123"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:123")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestUDPAddrFail(t *testing.T) {
	fv := flagvar.UDPAddr{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "999.999.999.999:1"})
	if err == nil {
		t.Fail()
	}
}

func TestUDPAddrs(t *testing.T) {
	fv := flagvar.UDPAddrs{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "127.0.0.1:123", "-udp-addr", ":80"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp", ":80")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUDPAddrsNetwork(t *testing.T) {
	fv := flagvar.UDPAddrs{Network: "udp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "127.0.0.1:123", "-udp-addr", ":80"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp4", ":80")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUDPAddrsFail(t *testing.T) {
	fv := flagvar.UDPAddrs{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addr", "")

	err := fs.Parse([]string{"-udp-addr", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestUDPAddrsCSV(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addrs-csv", "")

	err := fs.Parse([]string{"-udp-addrs-csv", "127.0.0.1:123,10.10.1.2:13"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp", "10.10.1.2:13")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUDPAddrsCSVNetwork(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{Network: "udp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addrs-csv", "")

	err := fs.Parse([]string{"-udp-addrs-csv", "127.0.0.1:123,10.10.1.2:13"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp4", "10.10.1.2:13")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUDPAddrsCSVSeparator(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addrs-csv", "")

	err := fs.Parse([]string{"-udp-addrs-csv", "127.0.0.1:123;:32"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp", ":32")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUDPAddrsCSVAccumulate(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addrs-csv", "")

	err := fs.Parse([]string{"-udp-addrs-csv", "127.0.0.1:123,10.10.1.2:432", "-udp-addrs-csv", "10.10.1.2:432"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUDPAddr("udp", "127.0.0.1:123")
	addr2, _ := net.ResolveUDPAddr("udp", "10.10.1.2:432")
	addr3, _ := net.ResolveUDPAddr("udp", "10.10.1.2:432")
	if !reflect.DeepEqual(fv.Values, []*net.UDPAddr{addr1, addr2, addr3}) {
		t.Fail()
	}
}

func TestUDPAddrsCSVFail(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "udp-addrs-csv", "")

	err := fs.Parse([]string{"-udp-addrs-csv", "999.999.999.1:1"})
	if err == nil {
		t.Fail()
	}
}
