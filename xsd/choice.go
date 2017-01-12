package xsd

// Choice http://www.w3schools.com/xml/el_choice.asp
type Choice struct {
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

// MaxOccurs returns maxOccurs
func (c Choice) MaxOccurs() string {
	return c.Max
}

// GetAllElements returns all inner Element's
func (c Choice) GetAllElements() []Element {
	elements := []Element{}
	if c.Elements != nil && len(c.Elements) > 0 {
		elements = append(elements, c.Elements...)
	}

	if c.Groups != nil {
		for _, g := range c.Groups {
			elements = append(elements, g.GetAllElements()...)
		}
	}

	if c.Choices != nil {
		for _, cc := range c.Choices {
			elements = append(elements, cc.GetAllElements()...)
		}
	}

	if c.Sequences != nil {
		for _, s := range c.Sequences {
			elements = append(elements, s.GetAllElements()...)
		}
	}

	return elements
}
