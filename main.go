package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

const (
	path       = `config/ignorer/cache`
	sourceRepo = `https://github.com/toptal/gitignore.git`
)

func main() {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      sourceRepo,
		Depth:    1,
		Progress: os.Stdout,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Fatal(err)
	}

	files := []gitignoreFile{
		{
			name: "Go",
			path: filepath.Join(`.`, path, `templates`, `Go.gitignore`),
		},
	}

	if err := readFile(os.Stdout, files...); err != nil {
		log.Fatal(err)
	}

	orders, err := readOrder(`./config/ignorer/cache/templates/order`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(orders)
}

type gitignoreFile struct {
	name string
	path string
}

func readFile(w io.Writer, files ...gitignoreFile) error {
	for _, file := range files {
		err := func(name, path string) error {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("read file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if _, err := io.WriteString(w, scanner.Text()+"\n"); err != nil {
					return fmt.Errorf("read file: %v", err)
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("read file: %v", err)
			}

			return nil
		}(file.name, file.path)

		if err != nil {
			return err
		}
	}

	return nil
}

func readOrder(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read order: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	orders := make(map[string]int)

	for n := 0; scanner.Scan(); {
		line := scanner.Text()
		if !isComment(line) {
			orders[line] = n
			n++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read order: %v", err)
	}

	return orders, nil
}

func isComment(line string) bool {
	return line != "" && line[0] == '#'
}

type template struct {
	name  string
	type_ string
}

type orderer struct {
	templates []template
	special   map[string]int
}

func (o *orderer) Len() int {
	return len(o.templates)
}

func (o *orderer) Swap(i, j int) {
	o.templates[i], o.templates[j] = o.templates[j], o.templates[i]
}

func (o *orderer) Less(i, j int) bool {
	for _, lessFn := range []func(int, int) bool{
		o.lessSpecial,
	} {
		less := lessFn
		switch {
		case less(i, j):
			return true
		case less(j, i):
			return false
		}
	}
	return o.lessName(i, j)
}

func (o *orderer) lessSpecial(i, j int) bool {
	in, jn := canon(o.templates[i].name), canon(o.templates[j].name)

	io, ok := o.special[in]
	if !ok {
		return false
	}

	jo, ok := o.special[jn]
	if !ok {
		return false
	}

	return io < jo
}

func (o *orderer) lessName(i, j int) bool {
	in, jn := canon(o.templates[i].name), canon(o.templates[j].name)
	return in < jn
}

func Sort(o orderer) orderer {
	sort.Sort(&o)
	return o
}

func canon(v string) string {
	return strings.ToLower(v)
}
