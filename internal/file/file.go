package file

import (
	"sort"
	"strings"
)

type File struct {
	Name string
	Typ  string
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

var typOrder = map[string]int{
	`.gitignore`: 0,
	`.patch`:     1,
	`.stack`:     2,
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
