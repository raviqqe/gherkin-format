package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/cucumber/gherkin/go/v27"
)

const featureFileExtension = ".feature"

func Format(r io.Reader, w io.Writer) error {
	d, err := gherkin.ParseGherkinDocument(r, func() string { return "" })

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, NewRenderer().Render(d))
	return err
}

func FormatPaths(ps []string) error {
	w := sync.WaitGroup{}
	es := make(chan error, len(ps))

	for _, p := range ps {
		w.Add(1)

		go func(p string) {
			defer w.Done()

			err := formatPath(p)

			if err != nil {
				es <- err
			}
		}(p)
	}

	w.Wait()

	if len(es) != 0 {
		return <-es
	}

	return nil
}

func formatPath(p string) error {
	s, err := os.Stat(p)

	if err != nil {
		return err
	} else if s.IsDir() {
		return FormatDirectory(p)
	}

	return FormatFile(p)
}

func FormatFile(s string) error {
	f, err := os.OpenFile(s, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	d, err := gherkin.ParseGherkinDocument(f, func() string { return s })

	if err != nil {
		return err
	}

	err = f.Truncate(0)

	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(f, NewRenderer().Render(d))
	return err
}

func FormatDirectory(d string) error {
	w := sync.WaitGroup{}
	es := make(chan error)

	err := filepath.Walk(d, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !i.IsDir() && filepath.Ext(p) == featureFileExtension {
			w.Add(1)

			go func() {
				defer w.Done()

				err := FormatFile(p)

				if err != nil {
					es <- err
				}
			}()
		}

		return nil
	})

	if err != nil {
		return err
	}

	w.Wait()

	if len(es) != 0 {
		return <-es
	}

	return nil
}
