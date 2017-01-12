package xsd

// Group http://www.w3schools.com/xml/el_group.asp
type Group struct {
	Annotation string     `xml:"annotation>documentation"`
	ID         string     `xml:"id,attr"`
	Min        string     `xml:"minOccurs,attr"`
	Max        string     `xml:"maxOccurs,attr"`
	Choices    []Choice   `xml:"choice"`
	Sequences  []Sequence `xml:"sequence"`
	All        []All      `xml:"all"`
}

func (g Group) MaxOccurs() string {
	return g.Max
}

// GetAllElements returns all internal elements
func (g Group) GetAllElements() []Element {
	elements := []Element{}
	if g.Choices != nil {
		for _, c := range g.Choices {
			elements = append(elements, c.GetAllElements()...)
		}
	}

	if g.Sequences != nil {
		for _, s := range g.Sequences {
			elements = append(elements, s.GetAllElements()...)
		}
	}

	if g.All != nil {
		for _, all := range g.All {
			elements = append(elements, all.GetAllElements()...)
		}
	}

	return elements
}
