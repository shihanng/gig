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

		if _, ok := filter[Canon(base)]; ok {
			files = append(files, File{Name: base, Typ: ext})
			filter[Canon(base)] = false
		}
	}

	undefineds := []string{}

	for k, v := range filter {
		if v {
			files = append(files, File{Name: k, Typ: ""})
			undefineds = append(undefineds, k)
		}
	}

	if len(undefineds) > 0 {
		return files, errors.Errorf("file: undefined template(s): %v", undefineds)
	}

	return files, nil
}

// Compose takes the contents of File from the given directory and join them together.
func Compose(w io.Writer, directory string, files ...File) error {
	for i, file := range files {
		err := func(name, ext string) error {
			var h string

			file, openErr := os.Open(filepath.Join(directory, name+ext))
			switch {
			case openErr == nil:
				h = header(name, ext)
			case os.IsNotExist(openErr):
				h = fmt.Sprintf("#!! ERROR: %s is undefined !!#\n", name)
			default:
				return errors.Wrap(openErr, "file: open file")
			}
			defer file.Close()

			if i > 0 {
				h = "\n" + h
			}

			if _, err := io.WriteString(w, h); err != nil {
				return errors.Wrap(err, "file: writing")
			}

			if openErr != nil {
				return nil
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

func Sort(files []File, special map[string]int) []File {
	specials := make([]File, 0, len(files))
	normals := make([]File, 0, len(files))

	for _, f := range files {
		if _, ok := special[Canon(f.Name)]; ok {
			specials = append(specials, f)
		} else {
			normals = append(normals, f)
		}
	}

	var specialFiles []File

	if len(specials) > 0 {
		specialSorter := sorter{
			files: specials,
			less: []lessFunc{
				lessSpecial(special),
				lessType,
			},
		}
		sort.Sort(&specialSorter)

		normals = append(normals, specialSorter.files[0])

		n := 1
		for ; n < len(specialSorter.files); n++ {
			if Canon(specialSorter.files[n].Name) != Canon(specialSorter.files[n-1].Name) {
				break
			}

			normals = append(normals, specialSorter.files[n])
		}

		specialFiles = specialSorter.files[n:]
	}

	normalSorter := sorter{
		files: normals,
		less: []lessFunc{
			lessName(special),
			lessType,
		},
	}
	sort.Sort(&normalSorter)

	return append(normalSorter.files, specialFiles...)
}

func Canon(v string) string {
	return strings.ToLower(v)
}

type lessFunc func(f, g File) bool

type sorter struct {
	files []File
	less  []lessFunc
}

func (s *sorter) Len() int {
	return len(s.files)
}

func (s *sorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

func (s *sorter) Less(i, j int) bool {
	p, q := s.files[i], s.files[j]

	var k int
	for k = 0; k < len(s.less)-1; k++ {
		less := s.less[k]

		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}

	return s.less[k](p, q)
}

func lessSpecial(special map[string]int) func(File, File) bool {
	return func(i, j File) bool {
		in, jn := Canon(i.Name), Canon(j.Name)

		io, ok := special[in]
		if !ok {
			return false
		}

		jo, ok := special[jn]
		if !ok {
			return false
		}

		return io < jo
	}
}

func lessName(special map[string]int) func(File, File) bool {
	return func(i, j File) bool {
		in, jn := Canon(i.Name), Canon(j.Name)

		_, iOK := special[in]
		_, jOK := special[jn]

		if iOK && jOK {
			return false
		}

		return in < jn
	}
}

func lessType(i, j File) bool {
	typOrder := map[string]int{
		`.gitignore`: 0,
		`.patch`:     1,
		`.stack`:     2,
	}

	in, jn := Canon(i.Name), Canon(j.Name)
	if in != jn {
		return false
	}

	it, jt := Canon(i.Typ), Canon(j.Typ)

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
