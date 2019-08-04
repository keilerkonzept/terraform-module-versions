package flagvar_test

import (
	"flag"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/sgreben/flagvar"
)

func TestAlternativeEither(t *testing.T) {
	uv := &flagvar.URL{}
	tv := &flagvar.Time{Layout: time.Kitchen}
	fv := flagvar.Alternative{
		Either: uv,
		Or:     tv,
	}
	fs := flag.FlagSet{}
	fs.Var(&fv, "url-or-time", "")

	u, _ := url.Parse("https://github.com/sgreben/flagvar")

	err := fs.Parse([]string{"-url-or-time", u.String()})
	if err != nil {
		t.Fail()
	}
	if !fv.EitherOk {
		t.Fail()
	}
	if !reflect.DeepEqual(uv.Value, u) {
		t.Fail()
	}
}

func TestAlternativeOr(t *testing.T) {
	uv := &flagvar.URL{}
	tv := &flagvar.Time{Layout: time.Kitchen}
	fv := flagvar.Alternative{
		Either: uv,
		Or:     tv,
	}
	fs := flag.FlagSet{}
	fs.Var(&fv, "url-or-time", "")

	u, _ := time.Parse(time.Kitchen, "10:30AM")

	err := fs.Parse([]string{"-url-or-time", u.Format(time.Kitchen)})
	if err != nil {
		t.Fail()
	}
	if fv.EitherOk {
		t.Fail()
	}
	if !reflect.DeepEqual(tv.Value, u) {
		t.Fail()
	}
}
