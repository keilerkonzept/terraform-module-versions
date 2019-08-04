package flagvar_test

import (
	"flag"
	"reflect"
	"testing"

	"github.com/sgreben/flagvar"
)

func TestAssignment(t *testing.T) {
	var fv flagvar.Assignment
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key=value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, struct {
		Key   string
		Value string
	}{Key: "key", Value: "value"}) {
		t.Fail()
	}
}

func TestAssignmentSeparator(t *testing.T) {
	fv := flagvar.Assignment{Separator: ":"}
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key:value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value, struct {
		Key   string
		Value string
	}{Key: "key", Value: "value"}) {
		t.Fail()
	}
}

func TestAssignmentFail(t *testing.T) {
	var fv flagvar.Assignment
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "keyXvalue"})
	if err == nil {
		t.Fail()
	}
}

func TestAssignments(t *testing.T) {
	var fv flagvar.Assignments
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key=value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []struct {
		Key   string
		Value string
	}{{Key: "key", Value: "value"}}) {
		t.Fail()
	}
}

func TestAssignmentsSeparator(t *testing.T) {
	fv := flagvar.Assignments{Separator: ":"}
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key:value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []struct {
		Key   string
		Value string
	}{{Key: "key", Value: "value"}}) {
		t.Fail()
	}
}

func TestAssignmentsFail(t *testing.T) {
	var fv flagvar.Assignments
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "keyXvalue"})
	if err == nil {
		t.Fail()
	}
}

func TestAssignmentsMap(t *testing.T) {
	var fv flagvar.AssignmentsMap
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key=value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, map[string]string{"key": "value"}) {
		t.Fail()
	}
}

func TestAssignmentsMapSeparator(t *testing.T) {
	fv := flagvar.AssignmentsMap{Separator: ":"}
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "key:value"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, map[string]string{"key": "value"}) {
		t.Fail()
	}
}

func TestAssignmentsMapFail(t *testing.T) {
	var fv flagvar.AssignmentsMap
	var fs flag.FlagSet
	fs.Var(&fv, "assignment", "")

	err := fs.Parse([]string{"-assignment", "keyXvalue"})
	if err == nil {
		t.Fail()
	}
}
