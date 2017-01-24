package xsd

import (
	"encoding/xml"
	"time"
)

// DefaultXSDDateFormat стандартныое представление типа xs:date
// https://www.w3.org/TR/xmlschema11-2/#date
const DefaultXSDDateFormat = "2006-01-02"

var defaultLocation *time.Location

func init() {
	defaultLocation, _ = time.LoadLocation("Europe/Moscow")
}

type Date struct {
	time.Time
}

// UnmarshalXMLAttr Позволяет правильно декодировать дату в формате "yyyy-mm-dd"
func (c *Date) UnmarshalXMLAttr(attr xml.Attr) error {
	t, err := time.ParseInLocation(time.RFC3339, attr.Value, defaultLocation)
	if err != nil {
		return err
	}

	c.Time = t
	return nil
}

// UnmarshalXML Позволяет правильно декодировать дату в формате "yyyy-mm-dd"
func (c *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.ParseInLocation(time.RFC3339, v, defaultLocation)
	if err != nil {
		return err
	}
	*c = Date{parse}
	return nil
}
