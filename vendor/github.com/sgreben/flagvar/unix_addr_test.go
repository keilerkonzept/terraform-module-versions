package flagvar_test

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestUnixAddr(t *testing.T) {
	fv := flagvar.UnixAddr{}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveUnixAddr("unix", "/example.sock")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestUnixAddrNetwork(t *testing.T) {
	fv := flagvar.UnixAddr{Network: "unixpacket"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveUnixAddr("unixpacket", "/example.sock")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestUnixAddrFail(t *testing.T) {
	fv := flagvar.UnixAddr{Network: "not-unix"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock"})
	if err == nil {
		t.Fail()
	}
}

func TestUnixAddrs(t *testing.T) {
	fv := flagvar.UnixAddrs{}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock", "-unix-addr", "/other.sock"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unix", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unix", "/other.sock")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUnixAddrsNetwork(t *testing.T) {
	fv := flagvar.UnixAddrs{Network: "unixpacket"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock", "-unix-addr", "/other.sock"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unixpacket", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unixpacket", "/other.sock")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUnixAddrsFail(t *testing.T) {
	fv := flagvar.UnixAddrs{Network: "not-unix"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addr", "")

	err := fs.Parse([]string{"-unix-addr", "/example.sock"})
	if err == nil {
		t.Fail()
	}
}

func TestUnixAddrsCSV(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addrs-csv", "")

	err := fs.Parse([]string{"-unix-addrs-csv", "/example.sock,/other.sock"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unix", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unix", "/other.sock")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUnixAddrsCSVNetwork(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{Network: "unixpacket"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addrs-csv", "")

	err := fs.Parse([]string{"-unix-addrs-csv", "/example.sock,/other.sock"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unixpacket", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unixpacket", "/other.sock")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUnixAddrsCSVSeparator(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addrs-csv", "")

	err := fs.Parse([]string{"-unix-addrs-csv", "/example.sock;:32"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unix", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unix", ":32")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestUnixAddrsCSVAccumulate(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addrs-csv", "")

	err := fs.Parse([]string{"-unix-addrs-csv", "/example.sock,/other.sock", "-unix-addrs-csv", "/other.sock"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveUnixAddr("unix", "/example.sock")
	addr2, _ := net.ResolveUnixAddr("unix", "/other.sock")
	addr3, _ := net.ResolveUnixAddr("unix", "/other.sock")
	if !reflect.DeepEqual(fv.Values, []*net.UnixAddr{addr1, addr2, addr3}) {
		t.Fail()
	}
}

func TestUnixAddrsCSVFail(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{Network: "not-unix"}
	var fs flag.FlagSet
	fs.Var(&fv, "unix-addrs-csv", "")

	err := fs.Parse([]string{"-unix-addrs-csv", "/example.sock"})
	if err == nil {
		t.Fail()
	}
}
