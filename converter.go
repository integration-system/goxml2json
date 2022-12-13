package xml2json

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type Converter struct {
	plugins []Plugin
}

func NewConverter(plugins ...Plugin) Converter {
	return Converter{
		plugins: plugins,
	}
}

// Convert converts the given XML document to JSON
func (s Converter) Convert(r io.Reader) (*bytes.Buffer, error) {
	root := &Node{}
	err := NewDecoder(r, s.plugins...).Decode(root)
	if err != nil {
		return nil, errors.WithMessage(err, "decode xml")
	}

	buf := new(bytes.Buffer)
	err = NewEncoder(buf, s.plugins...).Encode(root)
	if err != nil {
		return nil, errors.WithMessage(err, "encode json")
	}

	return buf, nil
}
