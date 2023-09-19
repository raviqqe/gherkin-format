package main

import (
	"fmt"
	"os"
)

func main() {
	if err := command(os.Args[1:]); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

func command(ss []string) error {
	args := getArguments(ss)
	s, err := os.Stat(args.Path)

	if err != nil {
		return err
	} else if s.IsDir() {
		return formatFiles(args.Path)
	}

	return formatFile(args.Path)
}
