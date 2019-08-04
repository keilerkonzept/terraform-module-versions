package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestInts(t *testing.T) {
	fv := flagvar.Ints{}
	var fs flag.FlagSet
	fs.Var(&fv, "ints", "")

	err := fs.Parse([]string{"-ints", "123", "-ints", "9090"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9090}) {
		t.Fail()
	}
}

func TestIntsBitSize(t *testing.T) {
	fv := flagvar.Ints{BitSize: 32}
	var fs flag.FlagSet
	fs.Var(&fv, "ints", "")

	err := fs.Parse([]string{"-ints", "123", "-ints", "9090"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9090}) {
		t.Fail()
	}
}

func TestIntsBase(t *testing.T) {
	fv := flagvar.Ints{BitSize: 32, Base: 16}
	var fs flag.FlagSet
	fs.Var(&fv, "ints", "")

	err := fs.Parse([]string{"-ints", "F", "-ints", "0"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{15, 0}) {
		t.Fail()
	}
}

func TestIntsFail(t *testing.T) {
	fv := flagvar.Ints{}
	var fs flag.FlagSet
	fs.Var(&fv, "ints", "")

	err := fs.Parse([]string{"-ints", "abc"})
	if err == nil {
		t.Fail()
	}
}

func TestIntsCSV(t *testing.T) {
	fv := flagvar.IntsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "ints-csv", "")

	err := fs.Parse([]string{"-ints-csv", "123,9492"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9492}) {
		t.Fail()
	}
}

func TestIntsCSVBitSize(t *testing.T) {
	fv := flagvar.IntsCSV{BitSize: 32}
	var fs flag.FlagSet
	fs.Var(&fv, "ints-csv", "")

	err := fs.Parse([]string{"-ints-csv", "123,9492"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9492}) {
		t.Fail()
	}
}

func TestIntsCSVSeparator(t *testing.T) {
	fv := flagvar.IntsCSV{Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "ints-csv", "")

	err := fs.Parse([]string{"-ints-csv", "123;9492"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9492}) {
		t.Fail()
	}
}

func TestIntsCSVAccumulate(t *testing.T) {
	fv := flagvar.IntsCSV{Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "ints-csv", "")

	err := fs.Parse([]string{"-ints-csv", "123,9492", "-ints-csv", "9492"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []int64{123, 9492, 9492}) {
		t.Fail()
	}
}

func TestIntsCSVFail(t *testing.T) {
	fv := flagvar.IntsCSV{}
	var fs flag.FlagSet
	fs.Var(&fv, "ints-csv", "")

	err := fs.Parse([]string{"-ints-csv", "third"})
	if err == nil {
		t.Fail()
	}
}
