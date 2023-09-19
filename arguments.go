package main

import (
	"github.com/docopt/docopt-go"
)

const usage = `Gherkin code formatter

Usage:
	gherkin2markdown <path>

Options:
	-h, --help  Show this help.`

type arguments struct {
	Path string `docopt:"<path>"`
}

func getArguments(ss []string) arguments {
	args := arguments{}
	parseArguments(usage, ss, &args)
	return args
}

func parseArguments(u string, ss []string, args interface{}) {
	opts, err := docopt.ParseArgs(u, ss, "0.1.0")

	if err != nil {
		panic(err)
	}

	if err := opts.Bind(args); err != nil {
		panic(err)
	}
}
