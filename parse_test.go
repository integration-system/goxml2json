package xml2json_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/integration-system/goxml2json"
	"github.com/stretchr/testify/suite"
)

func TestParse_Suite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &TestParse{})
}

type TestParse struct {
	suite.Suite
}

func (t *TestParse) SetupSuite() {}

type Product struct {
	ID       []int     `json:"id"`
	Price    []float64 `json:"price"`
	Deleted  []bool    `json:"deleted"`
	Nullable []any     `json:"nullable"`
}

type StringProduct struct {
	ID       []string `json:"id"`
	Price    []string `json:"price"`
	Deleted  []string `json:"deleted"`
	Nullable []string `json:"nullable"`
}

type MixedProduct struct {
	ID       []string  `json:"id"`
	Price    []float64 `json:"price"`
	Deleted  []string  `json:"deleted"`
	Nullable []string  `json:"nullable"`
}

const (
	productString = `
	<?xml version="1.0" encoding="UTF-8"?>	
		<id>42</id>
		<price>13.32</price>
		<deleted>true</deleted>
		<nullable>null</nullable>
		`
)

func (t *TestParse) TestAllJSTypeParsing() {
	converter := xml2json.NewConverter(
		xml2json.WithAttrPrefix("-"),
		xml2json.WithContentPrefix("#"),
		xml2json.WithTypeConverter(xml2json.Bool, xml2json.Int, xml2json.Float, xml2json.Null),
	)

	xml := strings.NewReader(productString)
	jsBuf, err := converter.Convert(xml)
	t.NoError(err, "could not parse test xml")

	product := Product{}
	err = json.Unmarshal(jsBuf.Bytes(), &product)
	t.NoError(err, "could not unmarshal test json")
	t.Equal(42, product.ID[0], "ID should match")
	t.Equal(13.32, product.Price[0], "price should match")
	t.Equal(true, product.Deleted[0], "deleted should match")
	t.Equal(nil, product.Nullable[0], "nullable should match")
}

func (t *TestParse) TestStringParsing() {
	converter := xml2json.NewConverter(
		xml2json.WithAttrPrefix("-"),
		xml2json.WithContentPrefix("#"),
	)

	xml := strings.NewReader(productString)
	jsBuf, err := converter.Convert(xml)
	t.NoError(err, "could not parse test xml")
	product := StringProduct{}
	err = json.Unmarshal(jsBuf.Bytes(), &product)
	t.NoError(err, "could not unmarshal test json")
	t.Equal("42", product.ID[0], "ID should match")
	t.Equal("13.32", product.Price[0], "price should match")
	t.Equal("true", product.Deleted[0], "deleted should match")
	t.Equal("null", product.Nullable[0], "nullable should match")
}

func (t *TestParse) TestMixedParsing() {
	converter := xml2json.NewConverter(
		xml2json.WithAttrPrefix("-"),
		xml2json.WithContentPrefix("#"),
		xml2json.WithTypeConverter(xml2json.Float),
	)

	xml := strings.NewReader(productString)
	jsBuf, err := converter.Convert(xml)
	t.NoError(err, "could not parse test xml")
	product := MixedProduct{}
	err = json.Unmarshal(jsBuf.Bytes(), &product)
	t.NoError(err, "could not unmarshal test json")
	t.Equal("42", product.ID[0], "ID should match")
	t.Equal(13.32, product.Price[0], "price should match")
	t.Equal("true", product.Deleted[0], "deleted should match")
	t.Equal("null", product.Nullable[0], "nullable should match")
}
