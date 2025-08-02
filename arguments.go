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
	args := Arguments{}
	s := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	s.BoolVar(&args.Check, "check", false, "check if files are formatted correctly")
	s.BoolVar(&args.Version, "version", false, "show version information")

	if err := s.Parse(ss); err != nil {
		return Arguments{}, err
	}

	args.Paths = s.Args()
	return args, nil
}
