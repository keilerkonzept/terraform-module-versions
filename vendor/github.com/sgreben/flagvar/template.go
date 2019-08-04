package flagvar

import (
	"fmt"
	"io/ioutil"
	"text/template"
)

// Template is a `flag.Value` for `text.Template` arguments.
// The value of the `Root` field is used as a root template when specified.
type Template struct {
	Root *template.Template

	Value *template.Template
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Template) Help() string {
	return "a go template"
}

// Set is flag.Value.Set
func (fv *Template) Set(v string) error {
	root := fv.Root
	if root == nil {
		root = template.New("")
	}
	t, err := root.New(fmt.Sprintf("%T(%p)", fv, fv)).Parse(v)
	if err == nil {
		fv.Value = t
	}
	return err
}

func (fv *Template) String() string {
	return fv.Text
}

// Templates is a `flag.Value` for `text.Template` arguments.
// The value of the `Root` field is used as a root template when specified.
type Templates struct {
	Root *template.Template

	Values []*template.Template
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *Templates) Help() string {
	return "a go template"
}

// Set is flag.Value.Set
func (fv *Templates) Set(v string) error {
	root := fv.Root
	if root == nil {
		root = template.New("")
	}
	t, err := root.New(fmt.Sprintf("%T(%p)", fv, fv)).Parse(v)
	if err == nil {
		fv.Texts = append(fv.Texts, v)
		fv.Values = append(fv.Values, t)
	}
	return err
}

func (fv *Templates) String() string {
	return fmt.Sprint(fv.Texts)
}

// Template is a `flag.Value` for `text.Template` arguments.
// The value of the `Root` field is used as a root template when specified.
// The value specified on the command line is the path to the template
type TemplateFile struct {
	Root *template.Template

	Value *template.Template
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TemplateFile) Help() string {
	return "file system path to a go template"
}

// Set is flag.Value.Set
func (fv *TemplateFile) Set(v string) error {

	_template, err := ioutil.ReadFile(v)

	if err != nil {
		return err
	}

	root := fv.Root
	if root == nil {
		root = template.New("")
	}
	t, err := root.New(fmt.Sprintf("%T(%p)", fv, fv)).Parse(string(_template))
	if err == nil {
		fv.Text = v
		fv.Value = t
	}
	return err
}

func (fv *TemplateFile) String() string {
	return fv.Text
}

