package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/rpoletaev/parsexsd/xsd"
	"golang.org/x/tools/imports"
)

var (
	// Struct field generated from an element attribute
	attr = `{{ define "Attr" }}{{ printf "  %s " (lintTitle .Name) }}{{ printf "%s ` + "`xml:\\\"%s,attr\\\" json:\\\",omitempty\\\"`" + `" (lint .Type) .Name }}
{{ end }}`

	// Struct field generated from an element child element
	child = `{{ define "Child" }}{{ printf "  %s " (lintTitle .Name) }}{{ if .List }}[]{{ end }}{{ printf "%s ` + "`xml:\\\"%s,omitempty\\\" json:\\\",omitempty\\\"`" + `" (typeName (fieldType .)) .Name }}
{{ end }}`

	// Struct field generated from the character data of an element
	cdata = `{{ define "Cdata" }}{{ printf "%s %s ` + "`xml:\\\",chardata\\\" json:\\\",omitempty\\\"`" + `" (lintTitle .Name) (lint .Type) }}
{{ end }}`

	// Struct generated from a non-trivial element (with children and/or attributes)
	elem = `{{ printf "// %s is generated from an XSD element\ntype %s struct {\n" (typeName .Name) (typeName .Name) }}{{ range $a := .Attribs }}{{ template "Attr" $a }}{{ end }}{{ range $c := .Children }}{{ template "Child" $c }}{{ end }} {{ if .Cdata }}{{ template "Cdata" . }}{{ end }} }
`
)

var (
	// The initialism pairs are based on the commonInitialisms found in golang/lint
	// https://github.com/golang/lint/blob/4946cea8b6efd778dc31dc2dbeb919535e1b7529/lint.go#L698-L738
	//
	initialismPairs = []string{
		"Api", "API",
		"Ascii", "ASCII",
		"Cpu", "CPU",
		"Css", "CSS",
		"Dns", "DNS",
		"Eof", "EOF",
		"Guid", "GUID",
		"Html", "HTML",
		"Https", "HTTPS",
		"Http", "HTTP",
		"Id", "ID",
		"Ip", "IP",
		"Json", "JSON",
		"Lhs", "LHS",
		"Qps", "QPS",
		"Ram", "RAM",
		"Rhs", "RHS",
		"Rpc", "RPC",
		"Sla", "SLA",
		"Smtp", "SMTP",
		"Sql", "SQL",
		"Ssh", "SSH",
		"Tcp", "TCP",
		"Tls", "TLS",
		"Ttl", "TTL",
		"Udp", "UDP",
		"Uid", "UID",
		"Ui", "UI",
		"Uuid", "UUID",
		"Uri", "URI",
		"Url", "URL",
		"Utf8", "UTF8",
		"Vm", "VM",
		"Xml", "XML",
		"Xsrf", "XSRF",
		"Xss", "XSS",
	}

	initialisms  = strings.NewReplacer(initialismPairs...)
	builtinTypes = [20]string{"bool", "byte", "complex128", "complex64", "error", "float32", "float64", "int", "int16", "int32", "int64", "int8", "rune", "string", "uint", "uint16", "uint32", "uint64", "uint8", "uintptr"}
)

// Generator is responsible for generating Go structs based on a given XML
// schema tree.
type generator struct {
	pkg      string
	prefix   string
	exported bool
	types    map[string]struct{}
}

func (g generator) do(out io.Writer, roots []*xsd.XmlTree) error {
	g.types = make(map[string]struct{})

	tt, err := prepareTemplates(g.prefix, g.exported)
	if err != nil {
		return fmt.Errorf("could not prepare templates: %s", err)
	}

	var res bytes.Buffer

	if g.pkg != "" {
		fmt.Fprintf(&res, "// generated by goxsd; DO NOT EDIT\n\npackage %s\n\n", g.pkg)
	}

	fmt.Fprintf(&res, `import (
		"plugin"
	)`)

	for _, e := range roots {
		if err := g.execute(e, tt, &res); err != nil {
			myerr := err.Error()
			myerr += fmt.Sprintf("\n%+v\n", e)
			myerr += fmt.Sprintln(tt.Name())
			return fmt.Errorf(myerr)
		}
	}

	buf, err := imports.Process("", res.Bytes(), &imports.Options{
		Fragment:  true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	})
	if err != nil {
		io.Copy(out, bytes.NewBuffer(res.Bytes()))
		return err
	}

	if _, err := io.Copy(out, bytes.NewBuffer(buf)); err != nil {
		return err
	}

	return nil
}

func (g generator) execute(root *xsd.XmlTree, tt *template.Template, out io.Writer) error {
	if _, ok := g.types[root.Name]; ok {
		return nil
	}
	if err := tt.Execute(out, root); err != nil {
		return err
	}
	g.types[root.Name] = struct{}{}

	for _, e := range root.Children {
		if !primitiveType(e) {
			if err := g.execute(e, tt, out); err != nil {
				fmt.Printf("%+v", e)
				return err
			}
		}
	}

	return nil
}

func prepareTemplates(prefix string, exported bool) (*template.Template, error) {
	typeName := func(name string) string {
		if isTime(name) {
			return name
		}

		if isBuiltinType(name) {
			return name
		}
		if prefix != "" {
			name = prefix + strings.Title(name)
		}
		if exported {
			name = strings.Title(name)
		}

		return lint(name)
	}

	fmap := template.FuncMap{
		"lint":      lint,
		"lintTitle": lintTitle,
		"typeName":  typeName,
		"fieldType": fieldType,
	}

	tt := template.New("yyy").Funcs(fmap)
	if _, err := tt.Parse(attr); err != nil {
		return nil, err
	}
	if _, err := tt.Parse(cdata); err != nil {
		return nil, err
	}
	if _, err := tt.Parse(child); err != nil {
		return nil, err
	}
	if _, err := tt.Parse(elem); err != nil {
		return nil, err
	}
	return tt, nil
}

func isTime(typeName string) bool {
	return strings.HasPrefix(typeName, "time.") || strings.HasPrefix(typeName, "xsd.")
}

func isBuiltinType(typeName string) bool {
	clearType := strings.TrimPrefix(typeName, "[]")
	for _, t := range builtinTypes {
		if clearType == t {
			return true
		}
	}

	return false
}

// If this is a chardata field, the field type must point to a
// struct, even if the element type is a built-in primitive.
func fieldType(e *xsd.XmlTree) string {
	if e.Cdata {
		return e.Name
	}
	return e.Type
}

func primitiveType(e *xsd.XmlTree) bool {
	if e.Cdata {
		return false
	}

	if isBuiltinType(e.Type) {
		return true
	}
	return false
}

func lint(s string) string {
	return dashToCamel(squish(initialisms.Replace(s)))
}

func lintTitle(s string) string {
	return lint(strings.Title(s))
}

func squish(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func dashToCamel(name string) string {
	name = strings.Replace(name, "_", "-", -1)
	name = strings.Replace(name, ".", "-", -1)

	s := strings.Split(name, "-")
	if len(s) > 1 {
		for i := 1; i < len(s); i++ {
			s[i] = strings.Title(s[i])
		}
		return strings.Join(s, "")
	}
	return name
}
