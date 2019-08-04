package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"regexp"

	"github.com/sgreben/flagvar"
)

func TestRegexp(t *testing.T) {
	fv := flagvar.Regexp{}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-z]+"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, regexp.MustCompile("[a-z]+")) {
		t.Fail()
	}
}

func TestRegexpPOSIX(t *testing.T) {
	fv := flagvar.Regexp{POSIX: true}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-z]+"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, regexp.MustCompilePOSIX("[a-z]+")) {
		t.Fail()
	}
}

func TestRegexpFail(t *testing.T) {
	fv := flagvar.Regexp{}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestRegexps(t *testing.T) {
	fv := flagvar.Regexps{}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-z]+"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []*regexp.Regexp{regexp.MustCompile("[a-z]+")}) {
		t.Fail()
	}
}

func TestRegexpsPOSIX(t *testing.T) {
	fv := flagvar.Regexps{POSIX: true}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-z]+"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []*regexp.Regexp{regexp.MustCompilePOSIX("[a-z]+")}) {
		t.Fail()
	}
}

func TestRegexpsFail(t *testing.T) {
	fv := flagvar.Regexps{}
	var fs flag.FlagSet
	fs.Var(&fv, "regexp", "")

	err := fs.Parse([]string{"-regexp", "[a-"})
	if err == nil {
		t.Fail()
	}
}
