package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	funcMap          template.FuncMap
	templateInstance *template.Template

	debug = flag.Bool("debug", false, "enable debug messages")
)

func init() {
	log.SetPrefix("gotpl: ")
	flag.Parse()
	funcMap = template.FuncMap{
		"enc64": func(in string) string {
			return base64.StdEncoding.EncodeToString([]byte(in))
		},
	}
}

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values_in io.Reader, out io.Writer, tpl_files ...string) error {

	if *debug {
		log.Printf("tpl_files: %#v", tpl_files)
	}

	templateInstance = template.New(tpl_files[0]).Funcs(funcMap)
	tpl, err := templateInstance.ParseFiles(tpl_files...)
	if err != nil {
		return fmt.Errorf("error parsing template(s): %v", err)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, values_in)
	if err != nil {
		return fmt.Errorf("failed to read standard input: %v", err)
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		return fmt.Errorf("failed to parse standard input values: %v", err)
	}

	if *debug {
		log.Printf("\n%#v\n", values)
	}

	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("failed to parse standard input execute: %v", err)
	}
	return nil
}

func main() {

	shift := 1
	if *debug {
		shift = 2
	}

	err := ExecuteTemplates(os.Stdin, os.Stdout, os.Args[shift:]...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
