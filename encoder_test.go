package xml2json

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type bio struct {
	Firstname string
	Lastname  string
	Hobbies   []string
	Misc      map[string]string
}

// TestEncode ensures that encode outputs the expected JSON document.
func TestEncode(t *testing.T) {
	assert := assert.New(t)

	author := bio{
		Firstname: "Bastien",
		Lastname:  "Gysler",
		Hobbies:   []string{"DJ", "Running", "Tennis"},
		Misc: map[string]string{
			"lineSeparator": "\u2028",
			"Nationality":   "Swiss",
			"City":          "Zürich",
			"foo":           "",
			"bar":           "\"quoted text\"",
			"esc":           "escaped \\ sanitized",
			"r":             "\r return line",
			"default":       "< >",
			"runeError":     "\uFFFD",
		},
	}

	// Build document
	root := &Node{}
	root.AddChild("firstname", &Node{
		Data: author.Firstname,
	})
	root.AddChild("lastname", &Node{
		Data: author.Lastname,
	})

	for _, h := range author.Hobbies {
		root.AddChild("hobbies", &Node{
			Data: h,
		})
	}

	misc := &Node{}
	for k, v := range author.Misc {
		misc.AddChild(k, &Node{
			Data: v,
		})
	}
	root.AddChild("misc", misc)
	var enc *Encoder

	// Convert to JSON string
	buf := new(bytes.Buffer)
	enc = NewEncoder(buf)

	err := enc.Encode(nil)
	assert.NoError(err)

	attr := WithAttrPrefix("test")
	attr.AddToEncoder(enc)
	content := WithContentPrefix("test2")
	content.AddToEncoder(enc)

	err = enc.Encode(root)
	assert.NoError(err)

	// Build SimpleJSON
	expectedResultBytes := []byte(`
{
  "firstname": [
    "Bastien"
  ],
  "lastname": [
    "Gysler"
  ],
  "hobbies": [
    "DJ",
    "Running",
    "Tennis"
  ],
  "misc": [
    {
      "lineSeparator": ["\u2028"],
      "Nationality": ["Swiss"],
      "City": ["Zürich"],
      "foo": [""],
      "bar": ["\"quoted text\""],
      "esc": ["escaped \\ sanitized"],
      "r": ["\r return line"],
      "default": ["< >"],
      "runeError": ["\uFFFD"]
    }
  ]
}
`)
	expectedResult := make(map[string]any)
	err = json.Unmarshal(expectedResultBytes, &expectedResult)
	assert.NoError(err)

	actualResult := make(map[string]any)
	err = json.Unmarshal(buf.Bytes(), &actualResult)
	assert.NoError(err)

	assert.EqualValues(expectedResult, actualResult)
}
