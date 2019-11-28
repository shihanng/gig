package file

import (
	"sort"
	"strings"
)

type File struct {
	Name string
	Typ  string
}

var typOrder = map[string]int{
	`.gitignore`: 0,
	`.patch`:     1,
	`.stack`:     2,
}

type Orderer struct {
	Files   []File
	Special map[string]int
}

func (o *Orderer) Len() int {
	return len(o.Files)
}

func (o *Orderer) Swap(i, j int) {
	o.Files[i], o.Files[j] = o.Files[j], o.Files[i]
}

func (o *Orderer) Less(i, j int) bool {
	for _, lessFn := range []func(int, int) bool{
		o.lessSpecial,
		o.lessName,
	} {
		less := lessFn
		switch {
		case less(i, j):
			return true
		case less(j, i):
			return false
		}
	}
	return o.lessType(i, j)
}

func (o *Orderer) lessSpecial(i, j int) bool {
	in, jn := Canon(o.Files[i].Name), Canon(o.Files[j].Name)

	io, ok := o.Special[in]
	if !ok {
		return false
	}

	jo, ok := o.Special[jn]
	if !ok {
		return false
	}

	return io < jo
}

func (o *Orderer) lessName(i, j int) bool {
	in, jn := Canon(o.Files[i].Name), Canon(o.Files[j].Name)
	return in < jn
}

func (o *Orderer) lessType(i, j int) bool {
	it, jt := Canon(o.Files[i].Typ), Canon(o.Files[j].Typ)

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

func Sort(o Orderer) Orderer {
	sort.Sort(&o)
	return o
}

func Canon(v string) string {
	return strings.ToLower(v)
}
