package extract

import (
	"bytes"
	"encoding/json"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
)

func Deserialize(r io.Reader) (data interface{}, err error) {
	buf := &bytes.Buffer{}
	io.Copy(buf, r)
	bufBytes := buf.Bytes()

	if err = json.NewDecoder(buf).Decode(&data); err == nil {
		return data, nil
	}

	buf.Reset()
	buf.Write(bufBytes)
	if err = yaml.NewDecoder(buf).Decode(&data); err == nil {
		return data, nil
	}

	buf.Reset()
	buf.Write(bufBytes)
	if _, err = toml.DecodeReader(buf, &data); err == nil {
		return data, nil
	}

	return nil, err
}
