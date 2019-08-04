package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestJSON(t *testing.T) {
	fv := flagvar.JSON{}
	var fs flag.FlagSet
	fs.Var(&fv, "json", "")

	err := fs.Parse([]string{"-json", `[1,{"a":2}]`})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, []interface{}{float64(1), map[string]interface{}{"a": float64(2)}}) {
		t.Fail()
	}
}

func TestJSONFail(t *testing.T) {
	fv := flagvar.JSON{}
	var fs flag.FlagSet
	fs.Var(&fv, "json", "")

	err := fs.Parse([]string{"-json", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestJSONs(t *testing.T) {
	fv := flagvar.JSONs{}
	var fs flag.FlagSet
	fs.Var(&fv, "json", "")

	err := fs.Parse([]string{"-json", `[1,{"a":2}]`})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []interface{}{[]interface{}{float64(1), map[string]interface{}{"a": float64(2)}}}) {
		t.Fail()
	}
}

func TestJSONsValue(t *testing.T) {
	type example struct {
		A string
		B int
	}
	fv := flagvar.JSONs{
		Value: func() interface{} {
			return &example{}
		},
	}
	var fs flag.FlagSet
	fs.Var(&fv, "json", "")

	err := fs.Parse([]string{"-json", `{"A":"abc","B":123}`})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []interface{}{&example{A: "abc", B: 123}}) {
		t.Fail()
	}
}

func TestJSONsFail(t *testing.T) {
	fv := flagvar.JSONs{}
	var fs flag.FlagSet
	fs.Var(&fv, "json", "")

	err := fs.Parse([]string{"-json", "[a-"})
	if err == nil {
		t.Fail()
	}
}
