package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.1.0"

func main() {
	if err := Run(os.Args[1:], os.Stdout); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

func Run(ss []string, w io.Writer) error {
	args, err := GetArguments(ss)

	if err != nil {
		return err
	} else if args.Version {
		_, err := fmt.Fprintln(w, version)
		return err
	} else if len(args.Paths) == 0 {
		return Format(os.Stdin, os.Stdout)
	} else if args.Check {
		return CheckPaths(args.Paths)
	}

	return FormatPaths(args.Paths)
}
