package main

import (
	"flag"
	"strings"
)

var authToken string
var freq int
var entries entryFlags

func initArgs() {
	flag.IntVar(&freq, "f", 30, "ip change detection frequency in seconds")
	flag.Var(&entries, "e", "list of entries in the form of zone/record")
	flag.StringVar(&authToken, "t", "", "cloudflare api auth token")
	flag.Parse()
}

type entryFlags []string

func (e *entryFlags) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func (e *entryFlags) String() string {
	sb := strings.Builder{}
	sb.WriteRune('[')
	for i, s := range *e {
		sb.WriteString(s)
		if i < len(*e)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune(']')
	return sb.String()
}
