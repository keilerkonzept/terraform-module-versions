package flagvar_test

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/sgreben/flagvar"
)

func TestTime(t *testing.T) {
	fv := flagvar.Time{}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	err := fs.Parse([]string{"-time", "2018-05-18T23:28:03+02:00"})
	if err != nil {
		t.Fail()
	}
	tv, _ := time.Parse(time.RFC3339, "2018-05-18T23:28:03+02:00")
	if !reflect.DeepEqual(fv.Value, tv) {
		t.Fail()
	}
}

func TestTimeLayout(t *testing.T) {
	fv := flagvar.Time{Layout: time.Kitchen}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	err := fs.Parse([]string{"-time", "10:30AM"})
	if err != nil {
		t.Fail()
	}
	tv, _ := time.Parse(time.Kitchen, "10:30AM")
	if !reflect.DeepEqual(fv.Value, tv) {
		t.Fail()
	}
}

func TestTimeFail(t *testing.T) {
	fv := flagvar.Time{}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	err := fs.Parse([]string{"-time", "[a-"})
	if err == nil {
		t.Fail()
	}
}

func TestTimes(t *testing.T) {
	fv := flagvar.Times{}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	tv, _ := time.Parse(time.RFC3339, "2018-05-18T23:28:03+02:00")
	err := fs.Parse([]string{"-time", "2018-05-18T23:28:03+02:00"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Values, []time.Time{tv}) {
		t.Fail()
	}
}

func TestTimesLayout(t *testing.T) {
	fv := flagvar.Times{Layout: time.Kitchen}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	err := fs.Parse([]string{"-time", "10:30AM"})
	if err != nil {
		t.Fail()
	}
	tv, _ := time.Parse(time.Kitchen, "10:30AM")
	if !reflect.DeepEqual(fv.Values, []time.Time{tv}) {
		t.Fail()
	}
}

func TestTimesFail(t *testing.T) {
	fv := flagvar.Times{}
	var fs flag.FlagSet
	fs.Var(&fv, "time", "")

	err := fs.Parse([]string{"-time", "[a-"})
	if err == nil {
		t.Fail()
	}
}
