package xsd

import (
	"fmt"
	"strings"
)

type builder struct {
	schemas    []Schema
	complTypes map[string]ComplexType
	simplTypes map[string]SimpleType
}

// NewBuilder creates a new initialized builder populated with the given
// xsdSchema slice.
func NewBuilder(schemas []Schema) *builder {
	return &builder{
		schemas:    schemas,
		complTypes: make(map[string]ComplexType),
		simplTypes: make(map[string]SimpleType),
	}
}

type XmlTree struct {
	Name     string
	Type     string
	List     bool
	Cdata    bool
	Attribs  []xmlAttrib
	Children []*XmlTree
}

type xmlAttrib struct {
	Name string
	Type string
}

// buildXML generates and returns a tree of XmlTree objects based on a set of
// parsed XSD schemas.
func (b *builder) BuildXML() []*XmlTree {
	var roots []Element
	for _, s := range b.schemas {
		for _, e := range s.Elements {
			roots = append(roots, e)
		}
		for _, t := range s.ComplexTypes {
			b.complTypes[t.Name] = t
		}
		for _, t := range s.SimpleTypes {
			b.simplTypes[t.Name] = t
		}
	}

	var xelems []*XmlTree
	for _, e := range roots {
		xelems = append(xelems, b.BuildFromElement(e))
	}

	return xelems
}

// buildFromElement builds an XmlTree from an xsdElement, recursively
// traversing the XSD type information to build up an XML element hierarchy.
func (b *builder) BuildFromElement(e Element) *XmlTree {
	xelem := &XmlTree{Name: e.Name, Type: e.Name}
	println("elementName: ", e.Name)
	println("elementType: ", e.Type)
	if IsList(e) {
		xelem.List = true
	}

	if !e.IsInlineType() {
		switch t := b.findType(e.Type).(type) {
		case ComplexType:
			b.BuildFromComplexType(xelem, t)
		case SimpleType:
			b.BuildFromSimpleType(xelem, t)
		case string:
			xelem.Type = t
		}
		return xelem
	}

	if e.ComplexType != nil { // inline complex type
		b.BuildFromComplexType(xelem, *e.ComplexType)
		return xelem
	}

	if e.SimpleType != nil { // inline simple type
		b.BuildFromSimpleType(xelem, *e.SimpleType)
		return xelem
	}

	return xelem
}

// buildFromComplexType takes an XmlTree and an xsdComplexType, containing
// XSD type information for XmlTree enrichment.
func (b *builder) BuildFromComplexType(xelem *XmlTree, t ComplexType) {
	if t.Sequence != nil { // Does the element have children?
		for _, e := range t.Sequence.GetAllElements() {
			xelem.Children = append(xelem.Children, b.BuildFromElement(e))
		}
	}

	if t.All != nil {
		for _, e := range t.All.GetAllElements() {
			xelem.Children = append(xelem.Children, b.BuildFromElement(e))
		}
	}

	if t.Choice != nil {
		for _, e := range t.Choice.GetAllElements() {
			xelem.Children = append(xelem.Children, b.BuildFromElement(e))
		}
	}

	if t.Group != nil {
		for _, e := range t.Group.GetAllElements() {
			xelem.Children = append(xelem.Children, b.BuildFromElement(e))
		}
	}

	if t.Attributes != nil {
		b.BuildFromAttributes(xelem, t.Attributes)
	}

	if t.ComplexContent != nil {
		b.BuildFromComplexContent(xelem, *t.ComplexContent)
	}

	if t.SimpleContent != nil {
		b.BuildFromSimpleContent(xelem, *t.SimpleContent)
	}
}

// buildFromSimpleType assumes restriction child and fetches the base value,
// assuming that value is of a XSD built-in data type.
func (b *builder) BuildFromSimpleType(xelem *XmlTree, t SimpleType) {
	switch tp := b.findType(t.Restriction.Base).(type) {
	case string:
		xelem.Type = tp
	case SimpleType:
		b.BuildFromSimpleType(xelem, tp)
	case ComplexType:
		b.BuildFromComplexType(xelem, tp)
	}
}

func (b *builder) BuildFromComplexContent(xelem *XmlTree, c ComplexContent) {
	if c.Extension != nil {
		b.BuildFromExtension(xelem, c.Extension)
	}
}

// A simple content can refer to a text-only complex type
func (b *builder) BuildFromSimpleContent(xelem *XmlTree, c SimpleContent) {
	if c.Extension != nil {
		b.BuildFromExtension(xelem, c.Extension)
	}

	if c.Restriction != nil {
		b.BuildFromRestriction(xelem, c.Restriction)
	}
}

// buildFromExtension extends an existing type, simple or complex, with a
// sequence.
func (b *builder) BuildFromExtension(xelem *XmlTree, e *Extension) {
	switch t := b.findType(e.Base).(type) {
	case ComplexType:
		b.BuildFromComplexType(xelem, t)
	case SimpleType:
		b.BuildFromSimpleType(xelem, t)
		// If element is of simpleType and has attributes, it must collect
		// its value as chardata.
		if e.Attributes != nil {
			xelem.Cdata = true
		}
	default:
		xelem.Type = t.(string)
		// If element is of built-in type but has attributes, it must collect
		// its value as chardata.
		if e.Attributes != nil {
			xelem.Cdata = true
		}
	}

	if e.Sequence != nil {
		for _, e := range e.Sequence {
			xelem.Children = append(xelem.Children, b.BuildFromElement(e))
		}
	}

	if e.Attributes != nil {
		b.BuildFromAttributes(xelem, e.Attributes)
	}
}

func (b *builder) BuildFromRestriction(xelem *XmlTree, r *Restriction) {
	switch t := b.findType(r.Base).(type) {
	case SimpleType:
		b.BuildFromSimpleType(xelem, t)
	case ComplexType:
		b.BuildFromComplexType(xelem, t)
	case ComplexContent:
		panic("Restriction on complex content is not implemented")
	default:
		panic("Unexpected base type to restriction")
	}
}

func (b *builder) BuildFromAttributes(xelem *XmlTree, attrs []Attribute) {
	for _, a := range attrs {
		attr := xmlAttrib{Name: a.Name}
		switch t := b.findType(a.Type).(type) {
		case SimpleType:
			// Get type name from simpleType
			// If Restriction.Base is a simpleType or complexType, we panic
			attr.Type = b.findType(t.Restriction.Base).(string)
		case string:
			// If empty, then simpleType is present as content, but we ignore
			// that now
			attr.Type = t
		}
		xelem.Attribs = append(xelem.Attribs, attr)
	}
}

// findType takes a type name and checks if it is a registered XSD type
// (simple or complex), in which case that type is returned. If no such
// type can be found, the XSD specific primitive types are mapped to their
// Go correspondents. If no XSD type was found, the type name itself is
// returned.
func (b *builder) findType(name string) interface{} {
	name = stripNamespace(name)
	if t, ok := b.complTypes[name]; ok {
		println("has ComplexType with name ", name)
		return t
	}
	if t, ok := b.simplTypes[name]; ok {
		println("has SimpleType with name ", name)
		fmt.Printf("%v#\n", t)
		return t
	}

	switch name {
	case "boolean":
		return "bool"
	case "language", "Name", "token", "duration", "anyURI":
		return "string"
	case "long", "short", "integer", "int":
		return "int64"
	case "unsignedShort":
		return "uint16"
	case "decimal", "double":
		return "float64"
	case "dateTime":
		return "time.Time"
	case "date":
		return "xsd.Date"
	case "base64Binary":
		return "[]byte"
	case "positiveInteger":
		return "uint64"
	default:
		return name
	}
}

func stripNamespace(name string) string {
	if s := strings.Split(name, ":"); len(s) > 1 {
		return s[len(s)-1]
	}
	return name
}
