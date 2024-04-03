package xml2json_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	xml2json "github.com/txix-open/goxml2json"
)

func TestDecoder_Suite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &TestDecoder{})
}

type TestDecoder struct {
	source string
	suite.Suite
}

func (t *TestDecoder) SetupSuite() {
	t.source = `<?xml version="1.0" encoding="UTF-8"?>
  <osm version="0.6" generator="CGImap 0.0.2">
   <bounds minlat="54.0889580" minlon="12.2487570" maxlat="54.0913900" maxlon="12.2524800"/>
   <node id="298884269" lat="54.0901746" lon="12.2482632" user="SvenHRO" uid="46882" visible="true" version="1" changeset="676636" timestamp="2008-09-21T21:37:45Z"/>
   <node id="261728686" lat="54.0906309" lon="12.2441924" user="PikoWinter" uid="36744" visible="true" version="1" changeset="323878" timestamp="2008-05-03T13:39:23Z"/>
   <node id="1831881213" version="1" changeset="12370172" lat="54.0900666" lon="12.2539381" user="lafkor" uid="75625" visible="true" timestamp="2012-07-20T09:43:19Z">
    <tag k="name" v="Neu Broderstorf"/>
    <tag k="traffic_sign" v="city_limit"/>
   </node>
   <foo>bar</foo>
  </osm>`
}

// TestDecode ensures that decode does not return any errors (not that useful)
func (t *TestDecoder) TestDecode() {
	// Decode XML document
	root := &xml2json.Node{}
	dec := xml2json.NewDecoder(strings.NewReader(t.source))
	err := dec.Decode(root)
	t.NoError(err)

	dec.SetAttributePrefix("test")
	dec.SetContentPrefix("test2")
	err = dec.DecodeWithCustomPrefixes(root, "test3", "test4")
	t.NoError(err)
}

func (t *TestDecoder) TestDecodeWithoutDefaultsAndExcludeAttributes() {
	// Decode XML document
	root := &xml2json.Node{}
	var err error
	var dec *xml2json.Decoder
	dec = xml2json.NewDecoder(
		strings.NewReader(t.source),
		xml2json.WithAttrPrefix(""),
		xml2json.ExcludeAttributes("version", "generator"),
	)
	err = dec.Decode(root)
	t.NoError(err)

	// Check that some attribute`s name has no prefix and has expected value
	t.Exactly(root.Children["osm"][0].Children["bounds"][0].Children["minlat"][0].Data, "54.0889580")
	// Check that some attributes are not present
	_, exists := root.Children["osm"][0].Children["version"]
	t.False(exists)
	_, exists = root.Children["osm"][0].Children["generator"]
	t.False(exists)
}

func (t *TestDecoder) TestTrim() {
	table := []struct {
		in       string
		expected string
	}{
		{in: "foo", expected: "foo"},
		{in: " foo", expected: "foo"},
		{in: "foo ", expected: "foo"},
		{in: " foo ", expected: "foo"},
		{in: "   foo   ", expected: "foo"},
		{in: "foo bar", expected: "foo bar"},
		{in: "\n\tfoo\n\t", expected: "foo"},
		{in: "\n\tfoo\n\tbar\n\t", expected: "foo\n\tbar"},
		{in: "", expected: ""},
		{in: "\n", expected: ""},
		{in: "\n\v", expected: ""},
		{in: "ending with ä", expected: "ending with ä"},
		{in: "ä and ä", expected: "ä and ä"},
	}

	for _, scenario := range table {
		got := xml2json.TrimNonGraphic(scenario.in)
		t.Equal(scenario.expected, got)
	}
}
