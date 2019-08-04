package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestFloat(t *testing.T) {
	fv := flagvar.Float{}
	var fs flag.FlagSet
	fs.Var(&fv, "float", "")

	err := fs.Parse([]string{"-float", "1.5"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, float64(1.5)) {
		t.Fail()
	}
}

func TestFloatHelp(t *testing.T) {
	fv := flagvar.Float{BitSize: 64}
	if !reflect.DeepEqual(fv.Help(), "a 64-bit float") {
		t.Fail()
	}
	fv = flagvar.Float{BitSize: 32}
	if !reflect.DeepEqual(fv.Help(), "a 32-bit float") {
		t.Fail()
	}
}

func TestFloatBitSize(t *testing.T) {
	fv := flagvar.Float{BitSize: 64}
	var fs flag.FlagSet
	fs.Var(&fv, "float", "")

	err := fs.Parse([]string{"-float", "1.234"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, float64(1.234)) {
		t.Fail()
	}
}

func TestFloatFail(t *testing.T) {
	fv := flagvar.Float{}
	var fs flag.FlagSet
	fs.Var(&fv, "float", "")

	err := fs.Parse([]string{"-float", "abc"})
	if err == nil {
		t.Fail()
	}
}
func TestFloats(t *testing.T) {
	fv := flagvar.Floats{}
	var fs flag.FlagSet
	fs.Var(&fv, "floats", "")

	err := fs.Parse([]string{"-floats", "1.5", "-floats", "0.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 0.0}) {
		t.Fail()
	}
}

func TestFloatsBitSize(t *testing.T) {
	fv := flagvar.Floats{BitSize: 32}
	var fs flag.FlagSet
	fs.Var(&fv, "floats", "")

	err := fs.Parse([]string{"-floats", "1.5", "-floats", "0.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 0.0}) {
		t.Fail()
	}
}

func TestFloatsFail(t *testing.T) {
	fv := flagvar.Floats{}
	var fs flag.FlagSet
	fs.Var(&fv, "floats", "")

	err := fs.Parse([]string{"-floats", "abc"})
	if err == nil {
		t.Fail()
	}
}

func TestFloatsCSV(t *testing.T) {
	fv := flagvar.FloatsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "floats-csv", "")

	err := fs.Parse([]string{"-floats-csv", "1.5,3.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 3.0}) {
		t.Fail()
	}
}

func TestFloatsCSVBitSize(t *testing.T) {
	fv := flagvar.FloatsCSV{BitSize: 32}
	var fs flag.FlagSet
	fs.Var(&fv, "floats-csv", "")

	err := fs.Parse([]string{"-floats-csv", "1.5,3.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 3.0}) {
		t.Fail()
	}
}

func TestFloatsCSVSeparator(t *testing.T) {
	fv := flagvar.FloatsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "floats-csv", "")

	err := fs.Parse([]string{"-floats-csv", "1.5;3.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 3.0}) {
		t.Fail()
	}
}

func TestFloatsCSVAccumulate(t *testing.T) {
	fv := flagvar.FloatsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "floats-csv", "")

	err := fs.Parse([]string{"-floats-csv", "1.5,3.0", "-floats-csv", "3.0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []float64{1.5, 3.0, 3.0}) {
		t.Fail()
	}
}

func TestFloatsCSVFail(t *testing.T) {
	fv := flagvar.FloatsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "floats-csv", "")

	err := fs.Parse([]string{"-floats-csv", "third"})
	if err == nil {
		t.Fail()
	}
}
