package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/text/encoding/charmap"

	"strings"

	"github.com/rpoletaev/parsexsd/xsd"
)

var (
	parsedFiles = make(map[string]struct{})

	repository, pckg, prefix string
	exported                 bool

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
	flag.StringVar(&repository, "o", "/go/src/bitbucket.org/losaped/goszakupki/repository", "Name of output file")
	flag.StringVar(&pckg, "p", "main", "Name of the Go package")
	flag.StringVar(&prefix, "x", "", "Name of the Go package")
	flag.BoolVar(&exported, "e", true, "Generate exported structs")
	flag.Parse()

	xsdFile := flag.Args()[0] //"/home/roma/Загрузки/zakupki/scheme4.4/fcsExport.xsd"
	fdir, _ := filepath.Split(xsdFile)
	version, err := xsd.GetSchemaVersion(filepath.Join(fdir, "IntegrationTypes.xsd"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Version is: ", version)
	pluginDir := path.Join(repository, version.String())
	if _, err = os.Stat(pluginDir); err != nil {
		log.Println("Plugin directory is not exsist: ", err)
		if err = os.Mkdir(pluginDir, 0770); err != nil {
			log.Println("Error on create plugin dir: ", err)
		}
	}

	output := path.Join(pluginDir, "plugin.go")
	log.Println(output)
	out := os.Stdout

	if out, err = os.Create(output); err != nil {
		log.Errorln("Could not create or truncate output file:", output)
		os.Exit(1)
	}

	s, err := parseXSDFile(xsdFile)
	if err != nil {
		log.Fatal(err)
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

	compiler := NewPluginCompiler("export"+version.String(), output)
	if err = compiler.BuildPlugin(); err != nil {
		println(err.Error())
		log.Fatal(err)
	}
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

func parseXSDFile(fname string) ([]xsd.Schema, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("Не удалось открыть файл: %s\n%v", fname, err)
	}
	defer f.Close()

	var schema xsd.Schema
	d := xml.NewDecoder(f)
	d.CharsetReader = makeCharsetReader
	if err := d.Decode(&schema); err != nil {
		log.Println("Не удалось декодировать схему, ", fname)
		return nil, err
	}

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
