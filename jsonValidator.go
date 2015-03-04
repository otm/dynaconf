package main

import (
	"bytes"
	"encoding/json"
	"log"
)

type jsonValidator struct{}

func (v *jsonValidator) run(b []byte) (err error) {
	log.Printf("Validator: JSON - Running")

	var js map[string]interface{}

	err = json.Unmarshal(b, &js)
	if err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			line, col, _ := findLineCol(bytes.NewReader(b), serr.Offset)
			log.Printf("Illigal JSON at line %v column %v: %v", line, col, err)
		} else {
			log.Printf("Illigal JSON: %v", err)
		}
		return err
	}

	log.Printf("Validator: JSON - OK")
	return err
}
