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

func FormatPaths(paths []string) error {
	return visitPaths(paths, func(s string) error {
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
	})
}

func CheckPaths(paths []string) error {
	return visitPaths(paths, func(s string) error {
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
	})
}

func visitPaths(paths []string, visit func(string) error) error {
	w := sync.WaitGroup{}
	es := make(chan error)

	for _, p := range paths {
		if err := filepath.Walk(p, func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if i.IsDir() || filepath.Ext(p) != featureFileExtension {
				return nil
			}

			w.Add(1)
			go func() {
				defer w.Done()

				if err := visit(p); err != nil {
					es <- err
				}
			}()

			return nil
		}); err != nil {
			return err
		}
	}

	w.Wait()

	if len(es) != 0 {
		return <-es
	}

	return nil
}
