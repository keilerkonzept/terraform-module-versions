package flagvar_test

import (
	"flag"
	"net"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestTCPAddr(t *testing.T) {
	fv := flagvar.TCPAddr{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "127.0.0.1:123"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:123")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestTCPAddrNetwork(t *testing.T) {
	fv := flagvar.TCPAddr{Network: "tcp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "127.0.0.1:123"})
	if err != nil {
		t.Fail()
	}
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:123")
	if !reflect.DeepEqual(fv.Value, addr) {
		t.Fail()
	}
}

func TestTCPAddrFail(t *testing.T) {
	fv := flagvar.TCPAddr{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "999.999.999.999:1"})
	if err == nil {
		t.Fail()
	}
}

func TestTCPAddrs(t *testing.T) {
	fv := flagvar.TCPAddrs{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "127.0.0.1:123", "-tcp-addr", ":80"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp", ":80")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestTCPAddrsNetwork(t *testing.T) {
	fv := flagvar.TCPAddrs{Network: "tcp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "127.0.0.1:123", "-tcp-addr", ":80"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp4", ":80")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestTCPAddrsFail(t *testing.T) {
	fv := flagvar.TCPAddrs{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addr", "")

	err := fs.Parse([]string{"-tcp-addr", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestTCPAddrsCSV(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addrs-csv", "")

	err := fs.Parse([]string{"-tcp-addrs-csv", "127.0.0.1:123,10.10.1.2:13"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp", "10.10.1.2:13")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestTCPAddrsCSVNetwork(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{Network: "tcp4"}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addrs-csv", "")

	err := fs.Parse([]string{"-tcp-addrs-csv", "127.0.0.1:123,10.10.1.2:13"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp4", "10.10.1.2:13")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestTCPAddrsCSVSeparator(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addrs-csv", "")

	err := fs.Parse([]string{"-tcp-addrs-csv", "127.0.0.1:123;:32"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp", ":32")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2}) {
		t.Fail()
	}
}

func TestTCPAddrsCSVAccumulate(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addrs-csv", "")

	err := fs.Parse([]string{"-tcp-addrs-csv", "127.0.0.1:123,10.10.1.2:432", "-tcp-addrs-csv", "10.10.1.2:432"})
	if err != nil {
		t.Fail()
	}
	addr1, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:123")
	addr2, _ := net.ResolveTCPAddr("tcp", "10.10.1.2:432")
	addr3, _ := net.ResolveTCPAddr("tcp", "10.10.1.2:432")
	if !reflect.DeepEqual(fv.Values, []*net.TCPAddr{addr1, addr2, addr3}) {
		t.Fail()
	}
}

func TestTCPAddrsCSVFail(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "tcp-addrs-csv", "")

	err := fs.Parse([]string{"-tcp-addrs-csv", "999.999.999.1:1"})
	if err == nil {
		t.Fail()
	}
}
