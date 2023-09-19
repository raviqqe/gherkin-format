package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cucumber/gherkin-go/v19"
)

const featureFileExtension = ".feature"

func formatFile(s string) error {
	f, err := os.OpenFile(s, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	d, err := gherkin.ParseGherkinDocument(f, func() string { return s })

	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	_, err = fmt.Fprint(f, newRenderer().Render(d))
	return err
}

func formatFiles(d string) error {
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

				err := formatFile(p)

				if err != nil {
					es <- err
					return
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
