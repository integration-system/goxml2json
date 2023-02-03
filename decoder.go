package xml2json

import (
	"encoding/xml"
	"fmt"
	"io"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
)

// A Decoder reads and decodes XML objects from an input stream.
type Decoder struct {
	reader          io.Reader
	err             error
	attributePrefix string
	contentPrefix   string
	excludeAttrs    map[string]bool
}

type element struct {
	parent *element
	n      *Node
	label  string
}

func (dec *Decoder) SetAttributePrefix(prefix string) {
	dec.attributePrefix = prefix
}

func (dec *Decoder) SetContentPrefix(prefix string) {
	dec.contentPrefix = prefix
}

func (dec *Decoder) ExcludeAttributes(attrs []string) {
	for _, attr := range attrs {
		dec.excludeAttrs[attr] = true
	}
}

func (dec *Decoder) DecodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
	dec.contentPrefix = contentPrefix
	dec.attributePrefix = attributePrefix
	return dec.Decode(root)
}

// NewDecoder returns a new decoder that reads from reader.
func NewDecoder(reader io.Reader, plugins ...Plugin) *Decoder {
	d := &Decoder{
		reader:          reader,
		contentPrefix:   "",
		attributePrefix: "",
		excludeAttrs:    make(map[string]bool),
	}
	for _, p := range plugins {
		d = p.AddToDecoder(d)
	}
	return d
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (dec *Decoder) Decode(root *Node) error {
	xmlDec := xml.NewDecoder(dec.reader)

	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	for {
		t, err := xmlDec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return errors.WithMessage(err, "xml decoder token")
		}

		switch se := t.(type) {
		case xml.StartElement:
			// Build new a new current element and link it to its parent
			elem = &element{
				parent: elem,
				n:      &Node{},
				label:  se.Name.Local,
			}

			// Extract attributes as children
			for _, a := range se.Attr {
				_, spaceFound := dec.excludeAttrs[a.Name.Space]
				_, localFound := dec.excludeAttrs[a.Name.Local]
				if spaceFound || localFound {
					continue
				}

				elem.n.AddChild(dec.attributePrefix+a.Name.Local, &Node{Data: a.Value})
			}
		case xml.CharData:
			// Extract XML data (if any)
			elem.n.Data = TrimNonGraphic(string(se))
		case xml.EndElement:
			// And add it to its parent list
			if elem.parent != nil {
				elem.parent.n.AddChild(elem.label, elem.n)
			}

			// Then change the current element to its parent
			elem = elem.parent
		}
	}

	dec.setPath("", root)

	return nil
}

func (dec *Decoder) setPath(path string, node *Node) {
	node.Label = path
	for label, nodes := range node.Children {
		childPath := label
		if path != "" {
			childPath = fmt.Sprintf("%s.%s", path, childPath)
		}

		for _, n := range nodes {
			dec.setPath(childPath, n)
		}
	}
}

// TrimNonGraphic returns a slice of the string s, with all leading and trailing
// non graphic characters and spaces removed.
//
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// and spaces, from categories L, M, N, P, S, Zs.
// Spacing characters are set by category Z and property Pattern_White_Space.
func TrimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i // copy i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	// If first is nil, it means there are no graphic characters
	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}
