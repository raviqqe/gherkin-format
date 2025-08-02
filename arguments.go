package main

import (
	"flag"
	"os"
)

type Arguments struct {
	Check   bool
	Paths   []string
	Version bool
}

func GetArguments(ss []string) (Arguments, error) {
	s := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c := s.Bool("check", false, "check if files are formatted correctly")
	v := s.Bool("version", false, "show version information")

	if err := s.Parse(ss); err != nil {
		return Arguments{}, err
	}

	return Arguments{
		Check:   *c,
		Paths:   s.Args(),
		Version: *v,
	}, nil
}
