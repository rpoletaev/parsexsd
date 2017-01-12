package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rpoletaev/parsexsd/xsd"
)

const fname = "/home/roma/Загрузки/zakupki/scheme4.4/fcsExport.xsd"

var (
	parsedFiles = make(map[string]struct{})

	output, pckg, prefix string
	exported             bool

	usage = `Usage: goxsd [options] <xsd_file>

Options:
  -o <file>     Destination file [default: stdout]
  -p <package>  Package name [default: goxsd]
  -e            Generate exported structs [default: false]
  -x <prefix>   Struct name prefix [default: ""]

goxsd is a tool for generating XML decoding/encoding Go structs, according
to an XSD schema.
`
)

func main() {
	flag.StringVar(&output, "o", "/home/roma/projects/go/src/bitbucket.org/losaped/fillStructs/gen/gen.go", "Name of output file")
	flag.StringVar(&pckg, "p", "goxsd", "Name of the Go package")
	flag.StringVar(&prefix, "x", "", "Name of the Go package")
	flag.BoolVar(&exported, "e", false, "Generate exported structs")
	flag.Parse()

	// if len(flag.Args()) != 1 {
	// 	fmt.Println(usage)
	// 	os.Exit(1)
	// }
	// xsdFile := flag.Arg(0)

	s, err := parseXSDFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	out := os.Stdout
	if output != "" {
		if out, err = os.Create(output); err != nil {
			fmt.Println("Could not create or truncate output file:", output)
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
		fmt.Println("Code generation failed unexpectedly:", err.Error())
		os.Exit(1)
	}
}

func parseXSDFile(fname string) ([]xsd.Schema, error) {
	f, err := os.Open(fname)
	if err != nil {
		println(err)
	}
	defer f.Close()

	return parse(f, fname)
}

func parse(r io.Reader, fname string) ([]xsd.Schema, error) {
	var schema xsd.Schema

	d := xml.NewDecoder(r)
	// handle special character sets
	if err := d.Decode(&schema); err != nil {
		println(err)
	}

	//newFname := fname + ".txt"
	//ioutil.WriteFile(newFname, []byte(fmt.Sprintf("%#v\n", pretty.Formatter(schema))), 0644)
	schemas := []xsd.Schema{schema}
	dir, file := filepath.Split(fname)
	parsedFiles[file] = struct{}{}
	for _, imp := range schema.Imports {
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
