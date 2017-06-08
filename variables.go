package main

import (
	"fmt"
	"strings"
)

// Make Variables conform to the Flag Value interface
// See https://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/
// and the FlagSet example here https://golang.org/pkg/flag/
type Variables map[string]string

func (v *Variables) String() string {
	return fmt.Sprint(*v)
}

func (v *Variables) Set(value string) error {
	t := strings.Split(value, "=")
	if len(t) != 2 {
		return fmt.Errorf("failed to parse %v, proper format is 'key=value'", value)
	}
	(*v)[t[0]] = t[1]
	return nil
}
