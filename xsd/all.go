package xsd

//All http://www.w3schools.com/xml/el_all.asp
type All struct {
	Annotation string    `xml:"annotation>documentation"`
	ID         string    `xml:"id,attr"`
	Min        string    `xml:"minOccurs,attr"`
	Max        string    `xml:"maxOccurs,attr"`
	Elements   []Element `xml:"element"`
}

//MaxOccurs implements HasMaxOccurs interface
func (a All) MaxOccurs() string {
	return a.Max
}

// GetAllElements returns all internal elements
func (all All) GetAllElements() []Element {
	elements := []Element{}
	if all.Elements != nil {
		elements = append(elements, all.Elements...)
		for _, e := range all.Elements {
			elements = append(elements, e.GetAllElements()...)
		}
	}
	return elements
}
