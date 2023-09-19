package main

import (
	"github.com/docopt/docopt-go"
)

const usage = `Gherkin code formatter

Usage:
	gherkin2markdown [<path>]

Options:
	-h, --help  Show this help.`

type arguments struct {
	Path string `docopt:"<path>"`
}

func getArguments(ss []string) (arguments, error) {
	args := arguments{}
	err := parseArguments(usage, ss, &args)
	return args, err
}

func parseArguments(u string, ss []string, args interface{}) error {
	opts, err := docopt.ParseArgs(u, ss, "0.1.0")

	if err != nil {
		return err
	}

	return opts.Bind(args)
}
