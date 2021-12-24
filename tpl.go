package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	templateInstance *template.Template
)

func init() {
	funcMap := template.FuncMap{
		"enc64": func(in []byte) string {
			return base64.StdEncoding.EncodeToString(in)
		},
	}
	templateInstance = template.New("default")
	templateInstance = templateInstance.Funcs(funcMap)
}

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values_in io.Reader, out io.Writer, tpl_files ...string) error {
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
		return fmt.Errorf("failed to parse standard input: %v", err)
	}

	err = tpl.Execute(out, values)
	if err != nil {
		return fmt.Errorf("failed to parse standard input: %v", err)
	}
	return nil
}

func main() {
	err := ExecuteTemplates(os.Stdin, os.Stdout, os.Args[1:]...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
