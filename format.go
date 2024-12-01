package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/cucumber/gherkin/go/v27"
	"github.com/samber/lo"
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
	return visitPaths(paths, func(p string) error {
		f, err := os.OpenFile(p, os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		d, err := gherkin.ParseGherkinDocument(f, func() string { return p })
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
	return visitPaths(paths, func(p string) error {
		bs, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		d, err := gherkin.ParseGherkinDocument(bytes.NewReader(bs), func() string { return p })
		if err != nil {
			return err
		}

		if string(bs) != NewRenderer().Render(d) {
			return errors.New("file not formatted: p")
		}

		return nil
	})
}

func visitPaths(paths []string, visit func(string) error) error {
	w := sync.WaitGroup{}
	es := make(chan error, 64)

	for _, p := range paths {
		err := filepath.Walk(p, func(p string, i os.FileInfo, err error) error {
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

			w.Wait()

			return nil
		})
		if err != nil {
			return err
		}
	}

	go func() {
		w.Wait()
		close(es)
	}()

	if es := lo.ChannelToSlice(es); len(es) > 0 {
		return errors.Join(es...)
	}

	return nil
}
