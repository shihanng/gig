package template

import (
	"sort"
	"strings"
)

type Template struct {
	Name  string
	Type_ string
}

var typeOrder = map[string]int{
	`.gitignore`: 0,
	`.patch`:     1,
	`.stack`:     2,
}

type Orderer struct {
	Templates []Template
	Special   map[string]int
}

func (o *Orderer) Len() int {
	return len(o.Templates)
}

func (o *Orderer) Swap(i, j int) {
	o.Templates[i], o.Templates[j] = o.Templates[j], o.Templates[i]
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
	in, jn := Canon(o.Templates[i].Name), Canon(o.Templates[j].Name)

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
	in, jn := Canon(o.Templates[i].Name), Canon(o.Templates[j].Name)
	return in < jn
}

func (o *Orderer) lessType(i, j int) bool {
	it, jt := Canon(o.Templates[i].Type_), Canon(o.Templates[j].Type_)

	io, ok := typeOrder[it]
	if !ok {
		return false
	}

	jo, ok := typeOrder[jt]
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
