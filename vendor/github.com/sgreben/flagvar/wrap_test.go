package flagvar_test

import (
	"flag"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestWrapPointer(t *testing.T) {
	sv := flagvar.Strings{}
	var p flag.Value = &sv
	fv := flagvar.WrapPointer{
		Value: &p,
	}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-pointer", "")

	err := fs.Parse([]string{"-wrap-pointer", "abc"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(sv.Values, []string{"abc"}) {
		t.Fail()
	}
}

func TestWrapFunc(t *testing.T) {
	sv := &flagvar.Strings{}
	fv := flagvar.WrapFunc(func() flag.Value { return sv })
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-func", "")

	err := fs.Parse([]string{"-wrap-func", "abc"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(sv.Values, []string{"abc"}) {
		t.Fail()
	}
}

func TestWrap(t *testing.T) {
	updated := 0
	sv := &flagvar.Strings{}
	fv := flagvar.Wrap{
		Value: sv,
		Updated: func() {
			updated++
		},
	}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap", "")

	err := fs.Parse([]string{"-wrap", "abc", "-wrap", "def"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(sv.Values, []string{"abc", "def"}) {
		t.Fail()
	}
	if updated != 2 {
		t.Fail()
	}
}

func TestWrapCSV(t *testing.T) {
	sv := &flagvar.Strings{}
	fv := flagvar.WrapCSV{Value: sv}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-csv", "")

	err := fs.Parse([]string{"-wrap-csv", "abc,def", "-wrap-csv", "xyz"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(sv.Values, []string{"abc", "def", "xyz"}) {
		t.Fail()
	}
}

func TestWrapCSVFail(t *testing.T) {
	sv := &flagvar.IP{}
	fv := flagvar.WrapCSV{Value: sv}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-csv", "")

	err := fs.Parse([]string{"-wrap-csv", "127.0.0.1,def"})
	if err == nil {
		t.Fail()
	}
}

func TestWrapCSVSeparator(t *testing.T) {
	sv := &flagvar.Strings{}
	fv := flagvar.WrapCSV{
		Value:     sv,
		Separator: ";",
	}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-csv", "")

	err := fs.Parse([]string{"-wrap-csv", "abc;def", "-wrap-csv", "xyz"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(sv.Values, []string{"abc", "def", "xyz"}) {
		t.Fail()
	}
}

func TestWrapCSVUpdated(t *testing.T) {
	sv := &flagvar.Assignment{}
	current := map[string]string{}
	final := map[string]string{}
	fv := flagvar.WrapCSV{
		Value: sv,
		UpdatedOne: func() {
			current[sv.Value.Key] = sv.Value.Value
		},
		UpdatedAll: func() {
			final = current
			current = map[string]string{}
		},
		StringFunc: func() string {
			var out []string
			for k, v := range final {
				out = append(out, fmt.Sprintf("%s=%s", k, v))
			}
			sort.Strings(out)
			return strings.Join(out, ",")
		},
	}
	var fs flag.FlagSet
	fs.Var(&fv, "wrap-csv", "")

	err := fs.Parse([]string{"-wrap-csv", "xyz=abc", "-wrap-csv", "abc=def,def=xyz"})
	if err != nil {
		t.Fail()
	}
	fmt.Println(final)
	if !reflect.DeepEqual(final, map[string]string{"abc": "def", "def": "xyz"}) {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.String(), `abc=def,def=xyz`) {
		t.Fail()
	}
}
