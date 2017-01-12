package xsd

// Any http://www.w3schools.com/xml/el_any.asp
type Any struct {
	Annotation string `xml:"annotation>documentation"`
	ID         string `xml:"id,attr"`
	Min        string `xml:"minOccurs,attr"`
	Max        string `xml:"maxOccurs,attr"`
}

func (a Any) MaxOccurs() string {
	return a.Max
}
