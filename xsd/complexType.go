package xsd

// ComplexType http://www.w3schools.com/xml/el_complextype.asp
type ComplexType struct {
	Name           string          `xml:"name,attr"`
	Abstract       string          `xml:"abstract,attr"`
	Annotation     string          `xml:"annotation>documentation"`
	Sequence       *Sequence       `xml:"sequence"`
	Group          *Group          `xml:"group"`
	All            *All            `xml:"all"`
	Choice         *Choice         `xml:"choice"`
	Attributes     []Attribute     `xml:"attribute"`
	ComplexContent *ComplexContent `xml:"complexContent"`
	SimpleContent  *SimpleContent  `xml:"simpleContent"`
}

// GetAllElements returns all inner elements
func (ct ComplexType) GetAllElements() []Element {
	elements := []Element{}
	if ct.Sequence != nil {
		elements = append(elements, ct.Sequence.GetAllElements()...)
	}

	if ct.Group != nil {
		elements = append(elements, ct.Group.GetAllElements()...)
	}

	if ct.All != nil {
		elements = append(elements, ct.All.GetAllElements()...)
	}

	if ct.Choice != nil {
		elements = append(elements, ct.Choice.GetAllElements()...)
	}

	return elements
}
