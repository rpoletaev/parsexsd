package xsd

import (
	"encoding/xml"
	"strings"
)

// Schema is the root of our Go representation of an XSD schema.
// http://www.w3schools.com/xml/el_schema.asp
type Schema struct {
	XMLName      xml.Name
	Ns           string        `xml:"xmlns,attr"`
	Imports      []Import      `xml:"import"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
}

// Import http://www.w3schools.com/xml/el_import.asp
type Import struct {
	Location string `xml:"schemaLocation,attr"`
}

// NS parses the namespace from a value in the expected format
// http://host/namespace/v1 returns `namespace`
func (s Schema) NS() string {
	split := strings.Split(s.Ns, "/")
	if len(split) > 2 {
		return split[len(split)-2]
	}
	return ""
}

// HavingMaxOccurs represent types which contains MaxOccurs attribute
type HavingMaxOccurs interface {
	MaxOccurs() string
}

// IsList returns true if maxOccurs = 'unbounded'
func IsList(hmo HavingMaxOccurs) bool {
	return hmo.MaxOccurs() == "unbounded"
}

// ComplexContent http://www.w3schools.com/xml/el_complexcontent.asp
type ComplexContent struct {
	Extension   *Extension   `xml:"extension"`
	Restriction *Restriction `xml:"restriction"`
}

// SimpleContent http://www.w3schools.com/xml/el_simpleContent.asp
type SimpleContent struct {
	Extension   *Extension   `xml:"extension"`
	Restriction *Restriction `xml:"restriction"`
}

// Extension http://www.w3schools.com/xml/el_extension.asp
type Extension struct {
	Base       string      `xml:"base,attr"`
	Attributes []Attribute `xml:"attribute"`
	Sequence   []Element   `xml:"sequence>element"`
}

// Attribute http://www.w3schools.com/xml/el_attribute.asp
type Attribute struct {
	Name       string `xml:"name,attr"`
	Type       string `xml:"type,attr"`
	Use        string `xml:"use,attr"`
	Annotation string `xml:"annotation>documentation"`
}

// SimpleType http://www.w3schools.com/xml/el_simpletype.asp
type SimpleType struct {
	Name        string      `xml:"name,attr"`
	Annotation  string      `xml:"annotation>documentation"`
	Restriction Restriction `xml:"restriction"`
}

// Restriction http://www.w3schools.com/xml/el_restriction.asp
type Restriction struct {
	Base        string        `xml:"base,attr"`
	Pattern     Pattern       `xml:"pattern"`
	Enumeration []Enumeration `xml:"enumeration"`
}

// Pattern http://www.w3schools.com/xml/schema_elements_ref.asp
type Pattern struct {
	Value string `xml:"value,attr"`
}

// Enumeration http://www.w3schools.com/xml/schema_elements_ref.asp
type Enumeration struct {
	Value string `xml:"value,attr"`
}
