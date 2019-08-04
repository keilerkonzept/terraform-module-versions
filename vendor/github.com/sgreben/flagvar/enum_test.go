package flagvar_test

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func ExampleEnum() {
	mode := flagvar.Enum{Choices: []string{"fast", "exact"}}
	fruit := flagvar.Enums{Choices: []string{"apple", "pear"}}
	letters := flagvar.EnumsCSV{
		Choices:       []string{"a", "A", "b", "B"},
		CaseSensitive: true,
	}
	levels := flagvar.EnumSetCSV{
		Choices:    []string{"debug", "info", "warn", "error"},
		Accumulate: true,
	}

	var fs flag.FlagSet
	fs.Var(&mode, "mode", "set a mode")
	fs.Var(&fruit, "fruit", "add a fruit")
	fs.Var(&letters, "letters", "select some letters")
	fs.Var(&levels, "levels", "enable log levels")

	fs.Parse([]string{
		"-fruit", "Apple",
		"-fruit", "apple",
		"-fruit", "PEAR",
		"-letters", "b,A,a",
		"-levels", "debug",
		"-levels", "debug,info,error",
		"-mode", "fast",
		"-mode", "exact",
	})

	fmt.Println("fruit:", fruit.Values)
	fmt.Println("letters:", letters.Values)
	fmt.Println("levels:", levels.Values())
	fmt.Println("mode:", mode.Value)

	// Output:
	// fruit: [apple apple pear]
	// letters: [b A a]
	// levels: [debug error info]
	// mode: exact
}

func TestEnum(t *testing.T) {
	fv := flagvar.Enum{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "first", "-enum", "second"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, "second") {
		t.Fail()
	}
}

func TestEnumCaseSensitive(t *testing.T) {
	fv := flagvar.Enum{Choices: []string{"first", "FIRST"}, CaseSensitive: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "FIRST"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, "FIRST") {
		t.Fail()
	}
}

func TestEnumFail(t *testing.T) {
	fv := flagvar.Enum{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "third"})
	if err == nil {
		t.Fail()
	}
}

func TestEnums(t *testing.T) {
	fv := flagvar.Enums{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"first"}) {
		t.Fail()
	}
}

func TestEnumsCaseSensitive(t *testing.T) {
	fv := flagvar.Enums{Choices: []string{"first", "FIRST", "fIrSt"}, CaseSensitive: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "FIRST"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"FIRST"}) {
		t.Fail()
	}
}

func TestEnumsFail(t *testing.T) {

	fv := flagvar.Enums{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "third"})
	if err == nil {
		t.Fail()
	}
}

func TestEnumSet(t *testing.T) {
	fv := flagvar.EnumSet{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum-set", "")

	err := fs.Parse([]string{"-enum-set", "first", "-enum-set", "first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values(), []string{"first"}) {
		t.Fail()
	}
}

func TestEnumSetCaseSensitive(t *testing.T) {
	fv := flagvar.EnumSet{Choices: []string{"first", "FIRST", "fIrSt"}, CaseSensitive: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "FIRST", "-enum", "first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values(), []string{"FIRST", "first"}) {
		t.Fail()
	}
}

func TestEnumSetFail(t *testing.T) {

	fv := flagvar.EnumSet{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum-set", "")

	err := fs.Parse([]string{"-enum-set", "third"})
	if err == nil {
		t.Fail()
	}
}

func TestEnumsCSV(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enums-csv", "")

	err := fs.Parse([]string{"-enums-csv", "first,second"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"first", "second"}) {
		t.Fail()
	}
}

func TestEnumsCSVCaseSensitive(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"first", "FIRST", "fIrSt"}, CaseSensitive: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "FIRST,first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"FIRST", "first"}) {
		t.Fail()
	}
}

func TestEnumsCSVSeparator(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"first", "second"}, Separator: ";"}
	var fs flag.FlagSet
	fs.Var(&fv, "enums-csv", "")

	err := fs.Parse([]string{"-enums-csv", "first;second"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"first", "second"}) {
		t.Fail()
	}
}

func TestEnumsCSVAccumulate(t *testing.T) {
	fv := flagvar.EnumsCSV{Choices: []string{"first", "second"}, Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enums-csv", "")

	err := fs.Parse([]string{"-enums-csv", "first,second", "-enums-csv", "second"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []string{"first", "second", "second"}) {
		t.Fail()
	}
}

func TestEnumsCSVFail(t *testing.T) {

	fv := flagvar.EnumsCSV{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enums-csv", "")

	err := fs.Parse([]string{"-enums-csv", "third"})
	if err == nil {
		t.Fail()
	}
}

func TestEnumSetCSV(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"first", "second"}, Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum-set-csv", "")

	err := fs.Parse([]string{"-enum-set-csv", "first,second", "-enum-set-csv", "first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values(), []string{"first", "second"}) {
		t.Fail()
	}
}

func TestEnumSetCSVCaseSensitive(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"first", "FIRST", "fIrSt"}, CaseSensitive: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum", "")

	err := fs.Parse([]string{"-enum", "FIRST,first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values(), []string{"FIRST", "first"}) {
		t.Fail()
	}
}

func TestEnumSetCSVSeparator(t *testing.T) {
	fv := flagvar.EnumSetCSV{Choices: []string{"first", "second"}, Separator: ";", Accumulate: true}
	var fs flag.FlagSet
	fs.Var(&fv, "enum-set-csv", "")

	err := fs.Parse([]string{"-enum-set-csv", "first;second", "-enum-set-csv", "first"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values(), []string{"first", "second"}) {
		t.Fail()
	}
}

func TestEnumSetCSVFail(t *testing.T) {

	fv := flagvar.EnumSetCSV{Choices: []string{"first", "second"}}
	var fs flag.FlagSet
	fs.Var(&fv, "enum-set-csv", "")

	err := fs.Parse([]string{"-enum-set-csv", "third"})
	if err == nil {
		t.Fail()
	}
}
