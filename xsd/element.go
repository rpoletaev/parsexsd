package xsd

// Element http://www.w3schools.com/xml/el_element.asp
type Element struct {
	Name        string       `xml:"name,attr"`
	Type        string       `xml:"type,attr"`
	Default     string       `xml:"default,attr"`
	Min         string       `xml:"minOccurs,attr"`
	Max         string       `xml:"maxOccurs,attr"`
	Annotation  string       `xml:"annotation>documentation"`
	ComplexType *ComplexType `xml:"complexType"` // inline complex type
	SimpleType  *SimpleType  `xml:"simpleType"`  // inline simple type
}

func (e Element) IsInlineType() bool {
	return e.Type == ""
}

func (e Element) MaxOccurs() string {
	return e.Max
}

func (e Element) GetAllElements() []Element {
	elements := []Element{}
	if e.ComplexType != nil {
		elements = append(elements, e.ComplexType.GetAllElements()...)
	}
	return elements
}
