package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type validator interface {
	run([]byte) error
}

type configTemplate struct {
	template string
	output   string
	validate validator
	command  *exec.Cmd
}

func newConfigTemplate(s string) (*configTemplate, error) {
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid template string")
	}

	cfgTemplate := &configTemplate{template: parts[0], output: parts[1]}
	if len(parts) > 2 {
		for i := 2; i < len(parts); i++ {
			pair := strings.Split(parts[i], "=")
			if len(pair) != 2 {
				log.Fatal("Illigal key=value pair", parts[i])
			}
			switch pair[0] {
			case "validate":
				if pair[1] != "json" {
					log.Fatal("Illigal validator: ", pair[1])
				}
				cfgTemplate.validate = &jsonValidator{}
			case "exec":
				log.Printf("Executing command: %v\n", pair[1])

				cmdParts := strings.Split(pair[1], " ")
				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
				cfgTemplate.command = cmd
				log.Printf("Adding post command: %v\n", pair[1])
			default:
				log.Fatal("Unknown command: ", pair[i])
			}
		}
	}

	return cfgTemplate, nil
}

func createTemplates(templateStrings []string) ([]configTemplate, error) {
	var templates []configTemplate
	for _, template := range templateStrings {
		log.Printf("Template: %s\n", template)
		c, err := newConfigTemplate(template)
		if err != nil {
			return nil, err
		}
		templates = append(templates, *c)
	}

	return templates, nil
}

func (c *configTemplate) render(params map[string]interface{}) error {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		params[pair[0]] = pair[1]
	}

	log.Printf("Reading template: %v\n", c.template)
	ct, err := template.ParseFiles(c.template)
	if err != nil {
		log.Fatal("Unable to parse template:", c.template)
	}

	b := new(bytes.Buffer)
	ct.Execute(b, params)

	if c.validate != nil {
		err := c.validate.run(b.Bytes())
		if err != nil {
			log.Fatal("Aborting execution, validation failed")
		}
	}

	c.print(b)
	c.runPostCommand()
	return nil
}

func (c *configTemplate) print(r io.Reader) {
	var b bytes.Buffer

	if flgs.dry {
		io.Copy(os.Stdout, r)
		return
	}

	log.Printf("Writing output to: %v", c.output)
	f, err := os.Create(c.output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if flgs.log {
		logWriter := bufio.NewWriter(&b)
		mw := io.MultiWriter(f, logWriter)
		io.Copy(mw, r)
		logWriter.Flush()
	} else {
		io.Copy(f, r)
	}
	if flgs.log {
		log.Printf("%v OUTPUT %v\n%s\n%s", strings.Repeat("-", 10), strings.Repeat("-", 30), b.String(), strings.Repeat("-", 68))
	}

	log.Printf("Output written to: %v", c.output)
}

func (c *configTemplate) runPostCommand() {
	if c.command == nil {
		return
	}

	if flgs.dry {
		log.Printf("Mode dry-run: Skipping post command...\n")
		return
	}

	log.Printf("Running post command: %v %v\n", c.command.Path, c.command.Args[1:])
	b, err := c.command.CombinedOutput()
	log.Printf("%v OUTPUT %v\n%s\n%s", strings.Repeat("-", 10), strings.Repeat("-", 30), b, strings.Repeat("-", 68))
	if err != nil {
		log.Fatal("Error detected when running post command: ", err)
	}

}
