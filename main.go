package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/text/encoding/charmap"

	"strings"

	"github.com/rpoletaev/parsexsd/xsd"
)

var (
	parsedFiles = make(map[string]struct{})

	output, pckg, prefix string
	exported             bool

	usage = `Usage: parsexsd [options] <xsd_file>

Options:
  -o <file>     Destination file [default: stdout]
  -p <package>  Package name [default: main]
  -e            Generate exported structs [default: true]
  -x <prefix>   Struct name prefix [default: ""]

parsexsd is a tool for generating XML decoding/encoding Go structs, according
to an XSD schema.
`
)

func main() {
	flag.StringVar(&output, "o", "/home/roma/projects/go/src/bitbucket.org/losaped/fillStructs/gen/gen.go", "Name of output file")
	flag.StringVar(&pckg, "p", "main", "Name of the Go package")
	flag.StringVar(&prefix, "x", "", "Name of the Go package")
	flag.BoolVar(&exported, "e", true, "Generate exported structs")
	flag.Parse()

	xsdFile := flag.Args()[0] //"/home/roma/Загрузки/zakupki/scheme4.4/fcsExport.xsd"
	s, err := parseXSDFile(xsdFile)
	if err != nil {
		log.Fatal(err)
	}

	out := os.Stdout
	if output != "" {
		if out, err = os.Create(output); err != nil {
			log.Errorln("Could not create or truncate output file:", output)
			os.Exit(1)
		}
	}

	bldr := xsd.NewBuilder(s)

	gen := generator{
		pkg:      pckg,
		prefix:   prefix,
		exported: exported,
	}

	if err := gen.do(out, bldr.BuildXML()); err != nil {
		log.Errorln("Code generation failed unexpectedly:", err.Error())
		os.Exit(1)
	}

	ver := s[0].GetVersion()
	if out != os.Stdout {
		compiler := NewPluginCompiler("zakup_export_v"+ver.String(), output)
		if err = compiler.BuildPlugin(); err != nil {
			log.Fatal(err)
		}
	}
}

func parseXSDFile(fname string) ([]xsd.Schema, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parse(f, fname)
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	lc := strings.ToLower(charset)
	if lc == "windows-1252" {
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	}
	if lc == "windows-1251" {
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	return nil, fmt.Errorf("Unsuported charset: '%s'", charset)
}

func parse(r io.Reader, fname string) ([]xsd.Schema, error) {
	var doc xsd.Document
	//var schema xsd.Schema

	d := xml.NewDecoder(r)
	d.CharsetReader = makeCharsetReader
	if err := d.Decode(&doc); err != nil {
		return nil, err
	}
	if doc.Comment == "" {
		_, f := filepath.Split(fname)
		return nil, fmt.Errorf("Не указана схема для документа: %s", f)
	}

	doc.Schema.ExtractVersion(doc.Comment)
	println("Schema version: ", doc.Schema.GetVersion().String())
	schemas := []xsd.Schema{doc.Schema}
	dir, file := filepath.Split(fname)
	parsedFiles[file] = struct{}{}

	for _, imp := range doc.Schema.Imports {
		if _, ok := parsedFiles[imp.Location]; ok {
			continue
		}
		s, err := parseXSDFile(filepath.Join(dir, imp.Location))
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, s...)
	}
	return schemas, nil
}
