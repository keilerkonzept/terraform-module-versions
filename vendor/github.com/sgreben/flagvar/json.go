package flagvar

import (
	"encoding/json"
	"strings"
)

// JSON is a `flag.Value` for JSON arguments.
type JSON struct {
	Value interface{}
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *JSON) Help() string {
	return "a JSON value"
}

// Set is flag.Value.Set
func (fv *JSON) Set(v string) error {
	fv.Text = v
	if fv.Value == nil {
		return json.Unmarshal([]byte(v), &fv.Value)
	}
	return json.Unmarshal([]byte(v), fv.Value)
}

func (fv *JSON) String() string {
	return fv.Text
}

// JSONs is a `flag.Value` for JSON arguments. If non-nil, the `Value` field is used to generate template values.
type JSONs struct {
	Value  func() interface{}
	Values []interface{}
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *JSONs) Help() string {
	return "a JSON value"
}

// Set is flag.Value.Set
func (fv *JSONs) Set(v string) (err error) {
	var value interface{}
	if fv.Value != nil {
		value = fv.Value()
		err = json.Unmarshal([]byte(v), value)
	} else {
		err = json.Unmarshal([]byte(v), &value)
	}
	if err == nil {
		fv.Texts = append(fv.Texts, v)
		fv.Values = append(fv.Values, value)
	}
	return err
}

func (fv *JSONs) String() string {
	return strings.Join(fv.Texts, ",")
}
