package xsd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Schema is the root of our Go representation of an XSD schema.
// http://www.w3schools.com/xml/el_schema.asp
type Schema struct {
	XMLName      xml.Name
	Ns           string        `xml:"xmlns,attr"`
	Comment      string        `xml:",comment"`
	Imports      []Import      `xml:"import"`
	Elements     []Element     `xml:"element"`
	ComplexTypes []ComplexType `xml:"complexType"`
	SimpleTypes  []SimpleType  `xml:"simpleType"`
	Version      Version
}

//GetSchemaVersion parse file and returns version of xsd
func GetSchemaVersion(fname string) (Version, error) {
	//<!-- FCS INTEGRATION_TYPES Integration Scheme, version 4.4.0, create date 21.07.2014 -->
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить версию схемы: %s\n%v", fname, err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	// comment with version must be on second line
	reader.ReadString('\n')
	comment, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить версию схемы: %s\n%v", fname, err)
	}

	// println("comment is: ", comment)
	if strings.HasPrefix(comment, "<!--") {
		re := regexp.MustCompile(`version ([\d+\.?]+)`)
		m := re.FindStringSubmatch(comment)
		if len(m) != 2 {
			return nil, fmt.Errorf("Схема версии не указана")
		}

		splitMatch := strings.Split(m[1], ".")
		version := make([]int, len(splitMatch))
		for i, val := range splitMatch {
			intVal, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return nil, fmt.Errorf("Неправильный формат версии xsd")
			}
			version[i] = intVal
		}

		return version, nil
	}
	return nil, fmt.Errorf("Version could not be found")
}

// Version represents slice of version numbers from major to minor
type Version []int

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

// String implements of Stringer interface
func (v Version) String() string {
	res := ""
	if len(v) == 0 {
		return res
	}

	for i, val := range v {
		if i < len(v)-1 {
			res += strconv.Itoa(val) + "."
		}
	}

	res += strconv.Itoa(v[len(v)-1])
	return res
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
