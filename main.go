package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

func Run(ss []string) error {
	args, err := GetArguments(ss)

	if err != nil {
		return err
	} else if args.Path == "" {
		return Format(os.Stdin, os.Stdout)
	}

	s, err := os.Stat(args.Path)

	if err != nil {
		return err
	} else if s.IsDir() {
		return FormatFiles(args.Path)
	}

	return FormatFile(args.Path)
}
