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
	} else if len(args.Paths) == 0 {
		return Format(os.Stdin, os.Stdout)
	}

	s, err := os.Stat(args.Paths[0])

	if err != nil {
		return err
	} else if s.IsDir() {
		return FormatFiles(args.Paths[0])
	}

	return FormatFile(args.Paths[0])
}
