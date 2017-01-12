package xsd

// Sequence http://www.w3schools.com/xml/el_sequence.asp
type Sequence struct {
	Annotation string     `xml:"annotation>documentation"`
	ID         string     `xml:"id,attr"`
	Min        string     `xml:"minOccurs,attr"`
	Max        string     `xml:"maxOccurs,attr"`
	Elements   []Element  `xml:"element"`
	Groups     []Group    `xml:"group"`
	Choices    []Choice   `xml:"choice"`
	Sequences  []Sequence `xml:"sequence"`
	Anies      []Any      `xml:"any"`
}

func (s Sequence) MaxOccurs() string {
	return s.Max
}

// GetAllElements returns all internal elements
func (s Sequence) GetAllElements() []Element {
	elements := []Element{}
	elements = append(elements, s.Elements...)

	if s.Groups != nil {
		for _, g := range s.Groups {
			elements = append(elements, g.GetAllElements()...)
		}
	}

	if s.Choices != nil {
		for _, c := range s.Choices {
			elements = append(elements, c.GetAllElements()...)
		}
	}

	if s.Sequences != nil {
		for _, ss := range s.Sequences {
			elements = append(elements, ss.GetAllElements()...)
		}
	}
	return elements
}
