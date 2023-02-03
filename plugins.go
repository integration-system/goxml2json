package xml2json

import (
	"strings"
)

// Plugin is added to an encoder or/and to an decoder to allow custom functionality at runtime
type Plugin interface {
	AddToEncoder(*Encoder) *Encoder
	AddToDecoder(*Decoder) *Decoder
}

// encoderTypeConverter a type converter overides the default string sanitization for encoding json
type encoderTypeConverter interface {
	Convert(string) string
}

// customTypeConverter converts strings to JSON types using a best guess approach, only parses the JSON types given
// when initialized via WithTypeConverter
type customTypeConverter struct {
	parseTypes []JSType
}

// WithTypeConverter allows customized js type conversion behavior by passing in the desired JSTypes
func WithTypeConverter(ts ...JSType) Plugin {
	return &customTypeConverter{parseTypes: ts}
}

func (tc *customTypeConverter) parseAsString(t JSType) bool {
	if t == String {
		return true
	}
	for i := 0; i < len(tc.parseTypes); i++ {
		if tc.parseTypes[i] == t {
			return false
		}
	}
	return true
}

// AddToEncoder Adds the type converter to the encoder
func (tc *customTypeConverter) AddToEncoder(e *Encoder) *Encoder {
	e.tc = tc
	return e
}

func (tc *customTypeConverter) AddToDecoder(d *Decoder) *Decoder {
	return d
}

func (tc *customTypeConverter) Convert(s string) string {
	// remove quotes if they exists
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s = s[1 : len(s)-1]
	}
	jsType := Str2JSType(s)
	if tc.parseAsString(jsType) {
		// add the quotes removed at the start of this func
		s = `"` + s + `"`
	}
	return s
}

type attrPrefixer string

// WithAttrPrefix appends the given prefix to the json output of xml attribute fields to preserve namespaces
func WithAttrPrefix(prefix string) Plugin {
	ap := attrPrefixer(prefix)
	return &ap
}

func (a *attrPrefixer) AddToEncoder(e *Encoder) *Encoder {
	e.attributePrefix = string((*a))
	return e
}

func (a *attrPrefixer) AddToDecoder(d *Decoder) *Decoder {
	d.attributePrefix = string((*a))
	return d
}

type contentPrefixer string

// WithContentPrefix appends the given prefix to the json output of xml content fields to preserve namespaces
func WithContentPrefix(prefix string) Plugin {
	c := contentPrefixer(prefix)
	return &c
}

func (c *contentPrefixer) AddToEncoder(e *Encoder) *Encoder {
	e.contentPrefix = string((*c))
	return e
}

func (c *contentPrefixer) AddToDecoder(d *Decoder) *Decoder {
	d.contentPrefix = string((*c))
	return d
}

type excluder []string

// ExcludeAttributes excludes some xml attributes, for example, xmlns:xsi, xsi:noNamespaceSchemaLocation
func ExcludeAttributes(attrs []string) Plugin {
	ex := excluder(attrs)
	return &ex
}

func (ex *excluder) AddToEncoder(e *Encoder) *Encoder {
	return e
}

func (ex *excluder) AddToDecoder(d *Decoder) *Decoder {
	d.ExcludeAttributes([]string((*ex)))
	return d
}

type allToArray struct{}

func AllAttrToArray() Plugin {
	return allToArray{}
}

func (p allToArray) AddToEncoder(e *Encoder) *Encoder {
	e.allAttributeToArray = true
	return e
}

func (p allToArray) AddToDecoder(d *Decoder) *Decoder {
	return d
}

type attrToArray struct {
	attrList []string
}

func AttrToArray(attrList []string) Plugin {
	return attrToArray{
		attrList: attrList,
	}
}

func (p attrToArray) AddToEncoder(e *Encoder) *Encoder {
	attrIsArrayMap := make(map[string]bool)
	for _, s := range p.attrList {
		attrIsArrayMap[s] = true
	}

	e.attrIsAlwaysAnArray = attrIsArrayMap
	return e
}

func (p attrToArray) AddToDecoder(d *Decoder) *Decoder {
	return d
}
