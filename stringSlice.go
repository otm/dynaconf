package main

import "fmt"

type stringslice []string

func (s *stringslice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *stringslice) Set(value string) error {
	if value == "" {
		return fmt.Errorf("Empty strings are not allowed")
	}

	*s = append(*s, value)
	return nil
}
