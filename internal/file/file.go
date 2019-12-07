package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-multierror"
)

type File struct {
	Name string
	Typ  string
}

func List(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, errors.Wrap(err, "file: read directory for list")
	}

	var names []string

	for _, f := range files {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)

		if ext == ".gitignore" {
			names = append(names, base)
		}
	}

	return names, nil
}

type IgnoreFile struct {
	gitignore string
	patch     string
	stack     []string
}

func lookup(directory string, items []string) ([]string, map[string]IgnoreFile, error) {
	ignoreFiles := make(map[string]IgnoreFile)
	unique := make([]string, 0, len(items))

	for _, item := range items {
		if _, ok := ignoreFiles[Canon(item)]; ok {
			continue
		}

		ignoreFiles[Canon(item)] = IgnoreFile{}

		unique = append(unique, item)
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, nil, errors.Wrap(err, "file: read directory")
	}

	for _, f := range files {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		splitted := strings.Split(base, ".")

		if ignoreFile, ok := ignoreFiles[Canon(splitted[0])]; ok {
			switch Canon(ext) {
			case ".gitignore":
				ignoreFile.gitignore = filename
			case ".patch":
				ignoreFile.patch = filename
			case ".stack":
				ignoreFile.stack = append(ignoreFile.stack, filename)
			}

			ignoreFiles[Canon(splitted[0])] = ignoreFile
		}
	}

	return unique, ignoreFiles, nil
}

func Generate(w io.Writer, directory string, items ...string) error {
	uniqueItems, ignoreFiles, err := lookup(directory, items)
	if err != nil {
		return err
	}

	writer := writer{
		directory:  directory,
		duplicates: make(map[string]bool),
	}

	var errs *multierror.Error

	for _, item := range uniqueItems {
		ignoreFile := ignoreFiles[Canon(item)]

		if ignoreFile.gitignore == "" {
			if _, err := fmt.Fprintf(w, "\n#!! ERROR: %s is undefined !!#\n", item); err != nil {
				return errors.Wrap(err, "file: writing")
			}

			errs = multierror.Append(errs, errors.Errorf("file: %s is undefined", item))

			continue
		}

		if err := writer.Write(w, ignoreFile.gitignore); err != nil {
			return err
		}

		if err := writer.Write(w, ignoreFile.patch); err != nil {
			return err
		}

		if err := writer.Write(w, ignoreFile.stack...); err != nil {
			return err
		}
	}

	return errs.ErrorOrNil()
}

type writer struct {
	directory  string
	duplicates map[string]bool
}

func (w *writer) Write(out io.Writer, filenames ...string) error {
	for _, filename := range filenames {
		if filename == "" {
			continue
		}

		err := func(filename string) error {
			ext := filepath.Ext(filename)
			base := strings.TrimSuffix(filename, ext)

			if _, err := io.WriteString(out, header(base, ext)); err != nil {
				return errors.Wrap(err, "file: creating header")
			}

			file, err := os.Open(filepath.Join(w.directory, filename))
			if err != nil {
				return errors.Wrapf(err, "file: open file: %s", filename)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				content := strings.TrimSpace(scanner.Text())
				if content != "" && content[0] != '#' && w.duplicates[content] {
					continue
				}

				if _, err := fmt.Fprintln(out, content); err != nil {
					return errors.Wrap(err, "file: writing content")
				}
				w.duplicates[content] = true
			}

			if err := scanner.Err(); err != nil {
				return errors.Wrap(err, "file: scanning")
			}

			return nil
		}(filename)

		if err != nil {
			return err
		}
	}

	return nil
}

func header(name, typ string) string {
	switch Canon(typ) {
	case ".patch":
		typ = "Patch "
	case ".stack":
		typ = "Stack "
	default:
		typ = ""
	}

	return fmt.Sprintf("\n### %s %s###\n", name, typ)
}

func Canon(v string) string {
	return strings.ToLower(v)
}
