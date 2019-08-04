package flagvar_test

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/sgreben/flagvar"
)

func ExampleTemplate() {
	fv := flagvar.Template{
		Root: template.New("example").Funcs(template.FuncMap{
			"toUpper": func(s string) string {
				return strings.ToUpper(s)
			},
		}),
	}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	fs.Parse([]string{"-template", `{{ toUpper "hello, world!" }}`})
	fv.Value.Execute(os.Stdout, nil)

	// Output:
	// HELLO, WORLD!
}

func TestTemplate(t *testing.T) {
	fv := flagvar.Template{}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{.Abc}}"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value.Root.String(), template.Must(template.New("").Parse("{{.Abc}}")).Root.String()) {
		t.Fail()
	}
}

func TestTemplateRoot(t *testing.T) {
	fv := flagvar.Template{
		Root: template.New(""),
	}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{.Abc}}"})
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(fv.Value.Root.String(), template.Must(template.New("").Parse("{{.Abc}}")).Root.String()) {
		t.Fail()
	}
}

func TestTemplateFail(t *testing.T) {
	fv := flagvar.Template{}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{...}}"})
	if err == nil {
		t.Fail()
	}
}

func TestTemplates(t *testing.T) {
	fv := flagvar.Templates{}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{.Abc}}"})
	if err != nil {
		t.Fail()
	}
	for _, tmp := range fv.Values {
		if !reflect.DeepEqual(tmp.Root.String(), template.Must(template.New("").Parse("{{.Abc}}")).Root.String()) {
			t.Fail()
		}
	}
}

func TestTemplatesRoot(t *testing.T) {
	fv := flagvar.Templates{
		Root: template.New(""),
	}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{.Abc}}"})
	if err != nil {
		t.Fail()
	}
	for _, tmp := range fv.Values {
		if !reflect.DeepEqual(tmp.Root.String(), template.Must(template.New("").Parse("{{.Abc}}")).Root.String()) {
			t.Fail()
		}
	}
}

func TestTemplatesFail(t *testing.T) {
	fv := flagvar.Templates{}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "{{...}}"})
	if err == nil {
		t.Fail()
	}
}

func TestTemplateFile(t *testing.T) {
	fv := flagvar.TemplateFile{}
	var fs flag.FlagSet
	fs.Var(&fv, "template", "")

	err := fs.Parse([]string{"-template", "./noSuchFile.tpl"})
	if err == nil {
		t.Fail()
	}
}
