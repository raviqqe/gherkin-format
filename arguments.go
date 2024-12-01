package main

import (
	"flag"
	"os"
)

type Arguments struct {
	Check bool
	Paths []string
}

func GetArguments(ss []string) (Arguments, error) {
	s := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c := s.Bool("check", false, "Check if files are formatted correctly.")

	err := s.Parse(ss)

	if err != nil {
		return Arguments{}, err
	}

	return Arguments{Check: *c, Paths: s.Args()}, nil
}
