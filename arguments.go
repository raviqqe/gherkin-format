package main

import (
	"github.com/docopt/docopt-go"
)

const usage = `Gherkin code formatter

Usage:
	gherkin2markdown [<path>]

Options:
	-h, --help  Show this help.`

type Arguments struct {
	Path string `docopt:"<path>"`
}

func GetArguments(ss []string) (Arguments, error) {
	args := Arguments{}
	err := parseArguments(usage, ss, &args)
	return args, err
}

func parseArguments(u string, ss []string, args interface{}) error {
	opts, err := docopt.ParseArgs(u, ss, "")

	if err != nil {
		return err
	}

	return opts.Bind(args)
}
