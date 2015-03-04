package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	version = "0.0.1"
	author  = "Nils Lagerkvist"
)

type flags struct {
	version   bool
	params    stringslice
	templates stringslice
	validate  string
	dry       bool
	log       bool
}

var flgs = flags{}

func init() {
	flag.BoolVar(&flgs.version, "version", false, "Show version and exit")
	flag.Var(&flgs.templates, "template", "The template definition to use")
	flag.Var(&flgs.params, "param", "Parameters to use when generating configuration")
	flag.BoolVar(&flgs.dry, "dry-run", false, "Output result to stdout instead of files")
	flag.BoolVar(&flgs.log, "log-files", false, "Log files written to disk")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] -template <template:outfile[:actions]>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `   Actions can be one of the following:
      validate=json - validates output as JSON
      exec=<command>

`)
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	if flgs.version {
		versionHandler()
	}

	if len(flgs.templates) == 0 {
		flag.Usage()
	}

	params, err := readParams(flgs.params)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := createTemplates(flgs.templates)
	if err != nil {
		log.Fatal(err)
	}

	err = render(templates, params)
	if err != nil {
		log.Fatal(err)
	}
}

func render(templates []configTemplate, params map[string]interface{}) error {
	for _, template := range templates {
		err := template.render(params)
		if err != nil {
			return err
		}
	}

	return nil
}

func readParams(rawParams stringslice) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	for _, param := range rawParams {
		pair := strings.Split(param, "=")
		if len(pair) != 2 {
			return nil, fmt.Errorf("Parameters must be on the form 'key=value'")
		}
		params[pair[0]] = pair[1]
	}

	return params, nil
}

func versionHandler() {
	fmt.Println(os.Args[0], version)
	fmt.Print(`Copyright (C) 2015 Nils Lagerkvist
License MIT: The MIT License <http://http://opensource.org/licenses/MIT>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by Nils Lagerkvist <http://github.com/otm/dynaconf>
`)
	os.Exit(0)
}
