package flagvar_test

import (
	"flag"
	"net/url"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestURL(t *testing.T) {
	fv := flagvar.URL{}
	var fs flag.FlagSet
	fs.Var(&fv, "url", "")

	err := fs.Parse([]string{"-url", "https://github.com/sgreben/flagvar"})
	if err != nil {
		t.Fail()
	}
	uv, _ := url.Parse("https://github.com/sgreben/flagvar")
	if !reflect.DeepEqual(fv.Value, uv) {
		t.Fail()
	}
}

func TestURLFail(t *testing.T) {
	fv := flagvar.URL{}
	var fs flag.FlagSet
	fs.Var(&fv, "url", "")

	err := fs.Parse([]string{"-url", ":s[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestURLs(t *testing.T) {
	fv := flagvar.URLs{}
	var fs flag.FlagSet
	fs.Var(&fv, "url", "")

	uv, _ := url.Parse("https://github.com/sgreben/flagvar")
	err := fs.Parse([]string{"-url", "https://github.com/sgreben/flagvar"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []*url.URL{uv}) {
		t.Fail()
	}
}

func TestURLsFail(t *testing.T) {
	fv := flagvar.URLs{}
	var fs flag.FlagSet
	fs.Var(&fv, "url", "")

	err := fs.Parse([]string{"-url", ":s[a-"})
	if err == nil {
		t.Fail()
	}
}
