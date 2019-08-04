package flagvar_test

import (
	"flag"
	"testing"
	"time"

	"github.com/sgreben/flagvar"
)

func TestAlternativeEmptyHelp(t *testing.T) {
	fv := flagvar.Alternative{}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestAlternativeHelp(t *testing.T) {
	fv0 := flagvar.Enum{Choices: []string{"apple", "banana"}}
	fv1 := flagvar.Enum{Choices: []string{"kiwi", "pear"}}
	fv := flagvar.Alternative{
		Either: &fv0,
		Or:     &fv1,
	}
	help := fv.Help()
	if help != `either one of [apple banana], or one of [kiwi pear]` {
		t.Fail()
	}
}

func TestAssignmentHelp(t *testing.T) {
	fv := flagvar.Assignment{}
	help := fv.Help()
	if help != `a key/value pair KEY=VALUE` {
		t.Fail()
	}
}

func TestAssignmentSeparatorHelp(t *testing.T) {
	fv := flagvar.Assignment{Separator: ":"}
	help := fv.Help()
	if help != `a key/value pair KEY:VALUE` {
		t.Fail()
	}
}

func TestAssignmentsHelp(t *testing.T) {
	fv := flagvar.Assignments{}
	help := fv.Help()
	if help != `a key/value pair KEY=VALUE` {
		t.Fail()
	}
}

func TestAssignmentsSeparatorHelp(t *testing.T) {
	fv := flagvar.Assignments{Separator: ":"}
	help := fv.Help()
	if help != `a key/value pair KEY:VALUE` {
		t.Fail()
	}
}

func TestAssignmentsMapHelp(t *testing.T) {
	fv := flagvar.AssignmentsMap{}
	help := fv.Help()
	if help != `a key/value pair KEY=VALUE` {
		t.Fail()
	}
}

func TestAssignmentsMapSeparatorHelp(t *testing.T) {
	fv := flagvar.AssignmentsMap{Separator: ":"}
	help := fv.Help()
	if help != `a key/value pair KEY:VALUE` {
		t.Fail()
	}
}

func TestCIDRHelp(t *testing.T) {
	fv := flagvar.CIDR{}
	help := fv.Help()
	if help != `a CIDR notation IP address and prefix length` {
		t.Fail()
	}
}

func TestCIDRsHelp(t *testing.T) {
	fv := flagvar.CIDRs{}
	help := fv.Help()
	if help != `a CIDR notation IP address and prefix length` {
		t.Fail()
	}
}

func TestCIDRsCSVHelp(t *testing.T) {
	fv := flagvar.CIDRsCSV{}
	help := fv.Help()
	if help != `","-separated list of CIDR notation IP addresses/prefix lengths` {
		t.Fail()
	}
}

func TestEnumHelp(t *testing.T) {
	fv := flagvar.Enum{Choices: []string{"pig", "boar"}}
	help := fv.Help()
	if help != `one of [pig boar]` {
		t.Fail()
	}
}

func TestEnumCaseSensitiveHelp(t *testing.T) {
	fv := flagvar.Enum{Choices: []string{"pig", "boar"}, CaseSensitive: true}
	help := fv.Help()
	if help != `one of [pig boar] (case-sensitive)` {
		t.Fail()
	}
}

func TestEnumsHelp(t *testing.T) {
	fv := flagvar.Enums{Choices: []string{"pig", "boar"}}
	help := fv.Help()
	if help != `one of [pig boar]` {
		t.Fail()
	}
}

func TestEnumsCaseSensitiveHelp(t *testing.T) {
	fv := flagvar.Enums{Choices: []string{"pig", "boar"}, CaseSensitive: true}
	help := fv.Help()
	if help != `one of [pig boar] (case-sensitive)` {
		t.Fail()
	}
}

func TestEnumsCSVHelp(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"pig", "boar"}}
	help := fv.Help()
	if help != `","-separated list of values from [pig boar]` {
		t.Fail()
	}
}

func TestEnumsCSVCaseSensitiveHelp(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"pig", "boar"}, CaseSensitive: true}
	help := fv.Help()
	if help != `","-separated list of values from [pig boar] (case-sensitive)` {
		t.Fail()
	}
}

func TestEnumSetHelp(t *testing.T) {
	fv := flagvar.EnumSet{Choices: []string{"pig", "boar"}}
	help := fv.Help()
	if help != `one of [pig boar]` {
		t.Fail()
	}
}

func TestEnumSetCaseSensitiveHelp(t *testing.T) {
	fv := flagvar.EnumSet{Choices: []string{"pig", "boar"}, CaseSensitive: true}
	help := fv.Help()
	if help != `one of [pig boar] (case-sensitive)` {
		t.Fail()
	}
}

func TestEnumSetCSVHelp(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"pig", "boar"}}
	help := fv.Help()
	if help != `","-separated list of values from [pig boar]` {
		t.Fail()
	}
}

func TestEnumSetCSVCaseSensitiveHelp(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"pig", "boar"}, CaseSensitive: true}
	help := fv.Help()
	if help != `","-separated list of values from [pig boar] (case-sensitive)` {
		t.Fail()
	}
}

func TestEnumSetCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"pig", "boar"}, Separator: ";"}
	help := fv.Help()
	if help != `";"-separated list of values from [pig boar]` {
		t.Fail()
	}
}

func TestFloatsHelp(t *testing.T) {
	fv := flagvar.Floats{}
	help := fv.Help()
	if help != `a float` {
		t.Fail()
	}
}

func TestFloatsBitSizeHelp(t *testing.T) {
	fv := flagvar.Floats{BitSize: 32}
	help := fv.Help()
	if help != `a 32-bit float` {
		t.Fail()
	}
}

func TestFloatsCSVHelp(t *testing.T) {
	fv := flagvar.FloatsCSV{}
	help := fv.Help()
	if help != `","-separated list of floats` {
		t.Fail()
	}
}

func TestFloatsCSVBitSizeHelp(t *testing.T) {
	fv := flagvar.FloatsCSV{BitSize: 64}
	help := fv.Help()
	if help != `","-separated list of 64-bit floats` {
		t.Fail()
	}
}

func TestFloatsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.FloatsCSV{Separator: ";"}
	help := fv.Help()
	if help != `";"-separated list of floats` {
		t.Fail()
	}
}

func TestGlobHelp(t *testing.T) {
	fv := flagvar.Glob{}
	help := fv.Help()
	if help != `a glob expression with separator '/'` {
		t.Fail()
	}
}

func TestGlobSeparatorHelp(t *testing.T) {
	fv := flagvar.Glob{Separators: &[]rune{':'}}
	help := fv.Help()
	if help != `a glob expression with separator ':'` {
		t.Fail()
	}
}

func TestGlobSeparatorsHelp(t *testing.T) {
	fv := flagvar.Glob{Separators: &[]rune{':', '+'}}
	help := fv.Help()
	if help != `a glob expression with separators [':' '+']` {
		t.Fail()
	}
}

func TestGlobNoSeparatorHelp(t *testing.T) {
	fv := flagvar.Glob{Separators: &[]rune{}}
	help := fv.Help()
	if help != `a glob expression` {
		t.Fail()
	}
}

func TestGlobsHelp(t *testing.T) {
	fv := flagvar.Globs{}
	help := fv.Help()
	if help != `a glob expression with separator '/'` {
		t.Fail()
	}
}

func TestGlobsSeparatorHelp(t *testing.T) {
	fv := flagvar.Globs{Separators: &[]rune{':'}}
	help := fv.Help()
	if help != `a glob expression with separator ':'` {
		t.Fail()
	}
}

func TestGlobsSeparatorsHelp(t *testing.T) {
	fv := flagvar.Globs{Separators: &[]rune{':', '+'}}
	help := fv.Help()
	if help != `a glob expression with separators [':' '+']` {
		t.Fail()
	}
}

func TestGlobsNoSeparatorHelp(t *testing.T) {
	fv := flagvar.Globs{Separators: &[]rune{}}
	help := fv.Help()
	if help != `a glob expression` {
		t.Fail()
	}
}

func TestIntsHelp(t *testing.T) {
	fv := flagvar.Ints{}
	help := fv.Help()
	if help != `an integer` {
		t.Fail()
	}
}

func TestIntsBitSizeBaseHelp(t *testing.T) {
	fv := flagvar.Ints{BitSize: 32, Base: 16}
	help := fv.Help()
	if help != `a 32-bit base 16 integer` {
		t.Fail()
	}
}

func TestIntsCSVHelp(t *testing.T) {
	fv := flagvar.IntsCSV{}
	help := fv.Help()
	if help != `","-separated list of integers` {
		t.Fail()
	}
}

func TestIntsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.IntsCSV{Separator: ":"}
	help := fv.Help()
	if help != `":"-separated list of integers` {
		t.Fail()
	}
}

func TestIntsCSVBitSizeBaseHelp(t *testing.T) {
	fv := flagvar.IntsCSV{BitSize: 32, Base: 12}
	help := fv.Help()
	if help != `","-separated list of 32-bit base 12 integers` {
		t.Fail()
	}
}

func TestIPHelp(t *testing.T) {
	fv := flagvar.IP{}
	help := fv.Help()
	if help != `an IP address` {
		t.Fail()
	}
}

func TestIPsHelp(t *testing.T) {
	fv := flagvar.IPs{}
	help := fv.Help()
	if help != `an IP address` {
		t.Fail()
	}
}

func TestIPsCSVHelp(t *testing.T) {
	fv := flagvar.IPsCSV{}
	help := fv.Help()
	if help != `","-separated list of IP addresses` {
		t.Fail()
	}
}

func TestIPsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.IPsCSV{Separator: "#"}
	help := fv.Help()
	if help != `"#"-separated list of IP addresses` {
		t.Fail()
	}
}

func TestTCPAddrHelp(t *testing.T) {
	fv := flagvar.TCPAddr{}
	help := fv.Help()
	if help != `a TCP address` {
		t.Fail()
	}
}

func TestTCPAddrsHelp(t *testing.T) {
	fv := flagvar.TCPAddrs{}
	help := fv.Help()
	if help != `a TCP address` {
		t.Fail()
	}
}

func TestTCPAddrsCSVHelp(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{}
	help := fv.Help()
	if help != `","-separated list of TCP addresses` {
		t.Fail()
	}
}

func TestTCPAddrsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.TCPAddrsCSV{Separator: "#"}
	help := fv.Help()
	if help != `"#"-separated list of TCP addresses` {
		t.Fail()
	}
}

func TestUDPAddrHelp(t *testing.T) {
	fv := flagvar.UDPAddr{}
	help := fv.Help()
	if help != `a UDP address` {
		t.Fail()
	}
}

func TestUDPAddrsHelp(t *testing.T) {
	fv := flagvar.UDPAddrs{}
	help := fv.Help()
	if help != `a UDP address` {
		t.Fail()
	}
}

func TestUDPAddrsCSVHelp(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{}
	help := fv.Help()
	if help != `","-separated list of UDP addresses` {
		t.Fail()
	}
}

func TestUDPAddrsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.UDPAddrsCSV{Separator: "#"}
	help := fv.Help()
	if help != `"#"-separated list of UDP addresses` {
		t.Fail()
	}
}

func TestUnixAddrHelp(t *testing.T) {
	fv := flagvar.UnixAddr{}
	help := fv.Help()
	if help != `a UNIX domain socket address` {
		t.Fail()
	}
}

func TestUnixAddrsHelp(t *testing.T) {
	fv := flagvar.UnixAddrs{}
	help := fv.Help()
	if help != `a UNIX domain socket address` {
		t.Fail()
	}
}

func TestUnixAddrsCSVHelp(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{}
	help := fv.Help()
	if help != `","-separated list of UNIX domain socket addresses` {
		t.Fail()
	}
}

func TestUnixAddrsCSVSeparatorHelp(t *testing.T) {
	fv := flagvar.UnixAddrsCSV{Separator: "#"}
	help := fv.Help()
	if help != `"#"-separated list of UNIX domain socket addresses` {
		t.Fail()
	}
}

func TestJSONHelp(t *testing.T) {
	fv := flagvar.JSON{}
	help := fv.Help()
	if help != `a JSON value` {
		t.Fail()
	}
}

func TestJSONsHelp(t *testing.T) {
	fv := flagvar.JSONs{}
	help := fv.Help()
	if help != `a JSON value` {
		t.Fail()
	}
}

func TestRegexpHelp(t *testing.T) {
	fv := flagvar.Regexp{}
	help := fv.Help()
	if help != `a regular expression` {
		t.Fail()
	}
}

func TestRegexpPOSIXHelp(t *testing.T) {
	fv := flagvar.Regexp{POSIX: true}
	help := fv.Help()
	if help != `a POSIX regular expression` {
		t.Fail()
	}
}

func TestRegexpsHelp(t *testing.T) {
	fv := flagvar.Regexps{}
	help := fv.Help()
	if help != `a regular expression` {
		t.Fail()
	}
}

func TestRegexpsPOSIXHelp(t *testing.T) {
	fv := flagvar.Regexps{POSIX: true}
	help := fv.Help()
	if help != `a POSIX regular expression` {
		t.Fail()
	}
}

func TestTemplateHelp(t *testing.T) {
	fv := flagvar.Template{}
	help := fv.Help()
	if help != `a go template` {
		t.Fail()
	}
}

func TestTemplatesHelp(t *testing.T) {
	fv := flagvar.Templates{}
	help := fv.Help()
	if help != `a go template` {
		t.Fail()
	}
}

func TestTimeHelp(t *testing.T) {
	fv := flagvar.Time{}
	help := fv.Help()
	if help != `a time, e.g. 2006-01-02T15:04:05Z07:00` {
		t.Fail()
	}
}

func TestTimeLayoutHelp(t *testing.T) {
	fv := flagvar.Time{Layout: time.Kitchen}
	help := fv.Help()
	if help != `a time, e.g. 3:04PM` {
		t.Fail()
	}
}

func TestTimesHelp(t *testing.T) {
	fv := flagvar.Times{}
	help := fv.Help()
	if help != `a time, e.g. 2006-01-02T15:04:05Z07:00` {
		t.Fail()
	}
}

func TestTimesLayoutHelp(t *testing.T) {
	fv := flagvar.Times{Layout: time.Kitchen}
	help := fv.Help()
	if help != `a time, e.g. 3:04PM` {
		t.Fail()
	}
}

func TestURLHelp(t *testing.T) {
	fv := flagvar.URL{}
	help := fv.Help()
	if help != `a URL` {
		t.Fail()
	}
}

func TestURLsHelp(t *testing.T) {
	fv := flagvar.URLs{}
	help := fv.Help()
	if help != `a URL` {
		t.Fail()
	}
}

func TestWrapPointerHelp(t *testing.T) {
	fv0 := &flagvar.Enum{Choices: []string{"apple", "banana"}}
	fv1 := flag.Value(fv0)
	fv := flagvar.WrapPointer{Value: &fv1}
	help := fv.Help()
	if help != fv0.Help() {
		t.Fail()
	}
}

func TestWrapPointerEmptyHelp(t *testing.T) {
	var fv0 flag.Value
	fv := flagvar.WrapPointer{Value: &fv0}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapPointerEmpty2Help(t *testing.T) {
	fv0 := flag.Value(&flagvar.File{})
	fv := flagvar.WrapPointer{Value: &fv0}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapFuncHelp(t *testing.T) {
	fv0 := flagvar.Enum{Choices: []string{"apple", "banana"}}
	fv := flagvar.WrapFunc(func() flag.Value { return &fv0 })
	help := fv.Help()
	if help != fv0.Help() {
		t.Fail()
	}
}

func TestWrapFuncEmptyHelp(t *testing.T) {
	fv := flagvar.WrapFunc(nil)
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapFuncEmpty2Help(t *testing.T) {
	fv := flagvar.WrapFunc(func() flag.Value { return nil })
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapHelp(t *testing.T) {
	fv0 := flagvar.Enum{Choices: []string{"apple", "banana"}}
	fv := flagvar.Wrap{Value: &fv0}
	help := fv.Help()
	if help != fv0.Help() {
		t.Fail()
	}
}

func TestWrapEmptyHelp(t *testing.T) {
	fv := flagvar.Wrap{}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapEmpty2Help(t *testing.T) {
	fv := flagvar.Wrap{Value: &flagvar.File{}}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapCSVEmptyHelp(t *testing.T) {
	fv := flagvar.WrapCSV{}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapCSVEmpty2Help(t *testing.T) {
	fv := flagvar.WrapCSV{Value: &flagvar.File{}}
	help := fv.Help()
	if help != "" {
		t.Fail()
	}
}

func TestWrapCSVHelp(t *testing.T) {
	fv0 := flagvar.Enum{Choices: []string{"apple", "banana"}}
	fv := flagvar.WrapCSV{Value: &fv0}
	help := fv.Help()
	if help != `","-separated values, each value one of [apple banana]` {
		t.Fail()
	}
}
