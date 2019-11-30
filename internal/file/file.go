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
)

type File struct {
	Name string
	Typ  string
}

// Filter retrieves File from directory based on the content of the given filter.
func Filter(directory string, filter map[string]bool) ([]File, error) {
	files := []File{}

	fList, err := ioutil.ReadDir(directory)
	if err != nil {
		return files, errors.Wrap(err, "file: read directory")
	}

	for _, f := range fList {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)

		if filter[Canon(base)] {
			files = append(files, File{Name: base, Typ: ext})
		}
	}

	return files, nil
}

// Compose takes the contents of File from the given directory and join them together.
func Compose(w io.Writer, directory string, files ...File) error {
	for i, file := range files {
		err := func(name, ext string) error {
			file, err := os.Open(filepath.Join(directory, name+ext))
			if err != nil {
				return errors.Wrap(err, "file: open file")
			}
			defer file.Close()

			h := header(name, ext)
			if i > 0 {
				h = "\n" + h
			}

			if _, err := io.WriteString(w, h); err != nil {
				return errors.Wrap(err, "file: writing")
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if _, err := io.WriteString(w, scanner.Text()+"\n"); err != nil {
					return errors.Wrap(err, "file: writing")
				}
			}

			if err := scanner.Err(); err != nil {
				return errors.Wrap(err, "file: scanning")
			}

			return nil
		}(file.Name, file.Typ)

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

	return fmt.Sprintf("### %s %s###\n", name, typ)
}

func Sort(f []File, special map[string]int) []File {
	s := Sorter{
		Files:   f,
		Special: special,
	}
	sort.Sort(&s)

	return s.Files
}

func Canon(v string) string {
	return strings.ToLower(v)
}

type Sorter struct {
	Files   []File
	Special map[string]int
}

func (s *Sorter) Len() int {
	return len(s.Files)
}

func (s *Sorter) Swap(i, j int) {
	s.Files[i], s.Files[j] = s.Files[j], s.Files[i]
}

func (s *Sorter) Less(i, j int) bool {
	for _, lessFn := range []func(int, int) bool{
		s.lessSpecial,
		s.lessName,
	} {
		less := lessFn

		switch {
		case less(i, j):
			return true
		case less(j, i):
			return false
		}
	}

	return s.lessType(i, j)
}

func (s *Sorter) lessSpecial(i, j int) bool {
	in, jn := Canon(s.Files[i].Name), Canon(s.Files[j].Name)

	io, ok := s.Special[in]
	if !ok {
		return false
	}

	jo, ok := s.Special[jn]
	if !ok {
		return false
	}

	return io < jo
}

func (s *Sorter) lessName(i, j int) bool {
	in, jn := Canon(s.Files[i].Name), Canon(s.Files[j].Name)
	return in < jn
}

func (s *Sorter) lessType(i, j int) bool {
	typOrder := map[string]int{
		`.gitignore`: 0,
		`.patch`:     1,
		`.stack`:     2,
	}

	it, jt := Canon(s.Files[i].Typ), Canon(s.Files[j].Typ)

	io, ok := typOrder[it]
	if !ok {
		return false
	}

	jo, ok := typOrder[jt]
	if !ok {
		return false
	}

	return io < jo
}
