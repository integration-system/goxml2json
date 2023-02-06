package xml2json_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/integration-system/goxml2json"
	"github.com/stretchr/testify/suite"
)

func TestEncoder_Suite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &TestEncoder{})
}

type TestEncoder struct {
	suite.Suite
}

func (t *TestEncoder) SetupSuite() {}

// TestEncode ensures that encode outputs the expected JSON document.
func (t *TestEncoder) TestEncode() {
	type bio struct {
		Firstname string
		Lastname  string
		Hobbies   []string
		Misc      map[string]string
	}
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
	root := &xml2json.Node{}
	root.AddChild("firstname", &xml2json.Node{
		Data: author.Firstname,
	})
	root.AddChild("lastname", &xml2json.Node{
		Data: author.Lastname,
	})

	for _, h := range author.Hobbies {
		root.AddChild("hobbies", &xml2json.Node{
			Data: h,
		})
	}

	misc := &xml2json.Node{}
	for k, v := range author.Misc {
		misc.AddChild(k, &xml2json.Node{
			Data: v,
		})
	}
	root.AddChild("misc", misc)
	var enc *xml2json.Encoder

	// Convert to JSON string
	buf := new(bytes.Buffer)
	enc = xml2json.NewEncoder(buf)

	err := enc.Encode(nil)
	t.NoError(err)

	attr := xml2json.WithAttrPrefix("test")
	attr.AddToEncoder(enc)
	content := xml2json.WithContentPrefix("test2")
	content.AddToEncoder(enc)

	err = enc.Encode(root)
	t.NoError(err)

	// Build SimpleJSON
	expectedResultBytes := []byte(`
{
  "firstname": "Bastien",
  "lastname": "Gysler",
  "hobbies": [
    "DJ",
    "Running",
    "Tennis"
  ],
  "misc": {
      "lineSeparator": "\u2028",
      "Nationality": "Swiss",
      "City": "Zürich",
      "foo": "",
      "bar": "\"quoted text\"",
      "esc": "escaped \\ sanitized",
      "r": "\r return line",
      "default": "< >",
      "runeError": "\uFFFD"
  }
}
`)
	expectedResult := make(map[string]any)
	err = json.Unmarshal(expectedResultBytes, &expectedResult)
	t.NoError(err)

	actualResult := make(map[string]any)
	err = json.Unmarshal(buf.Bytes(), &actualResult)
	t.NoError(err)

	t.EqualValues(expectedResult, actualResult)
}
