package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/gobwas/glob"

	"github.com/sgreben/flagvar"
)

func TestGlob(t *testing.T) {
	fv := flagvar.Glob{}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "**.go"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, glob.MustCompile("**.go")) {
		t.Fail()
	}
}

func TestGlobNoSeparators(t *testing.T) {
	fv := flagvar.Glob{Separators: &[]rune{}}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "**.go"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, glob.MustCompile("**.go")) {
		t.Fail()
	}
}

func TestGlobSeparators(t *testing.T) {
	fv := flagvar.Glob{Separators: &[]rune{';'}}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "**.go"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, glob.MustCompile("**.go")) {
		t.Fail()
	}
}

func TestGlobFail(t *testing.T) {
	fv := flagvar.Glob{}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestGlobs(t *testing.T) {
	fv := flagvar.Globs{}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "**.go"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []glob.Glob{glob.MustCompile("**.go")}) {
		t.Fail()
	}
}

func TestGlobsFail(t *testing.T) {
	fv := flagvar.Globs{}
	var fs flag.FlagSet
	fs.Var(&fv, "glob", "")

	err := fs.Parse([]string{"-glob", "[a-"})
	if err == nil {
		t.Fail()
	}
}
