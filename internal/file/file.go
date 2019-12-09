package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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
			names = append(names, Canon(base))
		}
	}

	sort.Strings(names)

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
	ew := &errWriter{w: w}

	var errs *multierror.Error

	for _, item := range uniqueItems {
		ignoreFile := ignoreFiles[Canon(item)]

		if ignoreFile.gitignore == "" {
			ew.fprintf("\n#!! ERROR: %s is undefined !!#\n", item)

			errs = multierror.Append(errs, errors.Errorf("file: %s is undefined", item))

			continue
		}

		files := append([]string{ignoreFile.gitignore, ignoreFile.patch}, ignoreFile.stack...)
		if err := writer.Write(ew, files...); err != nil {
			return err
		}
	}

	if ew.err != nil {
		return ew.err
	}

	return errs.ErrorOrNil()
}

type writer struct {
	directory  string
	duplicates map[string]bool
}

func (w *writer) Write(out *errWriter, filenames ...string) error {
	var err error
	for _, filename := range filenames {
		if filename == "" || err != nil {
			continue
		}

		err = func(filename string) error {
			ext := filepath.Ext(filename)
			base := strings.TrimSuffix(filename, ext)

			out.fprintf(header(base, ext))

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

				out.fprintf("%s\n", content)
				w.duplicates[content] = true
			}

			return errors.Wrap(scanner.Err(), "file: scanning")
		}(filename)
	}

	return err
}

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) fprintf(format string, a ...interface{}) {
	if ew.err != nil {
		return
	}

	_, err := fmt.Fprintf(ew.w, format, a...)
	ew.err = errors.Wrap(err, "file: writing")
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
