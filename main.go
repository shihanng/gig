package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shihanng/gi/internal/order"
	"github.com/shihanng/gi/internal/template"
	"gopkg.in/src-d/go-git.v4"
)

const (
	path       = `config/gi/cache`
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

	args := os.Args[1:]
	languages := make(map[string]bool, len(args))

	for _, arg := range args {
		languages[template.Canon(arg)] = true
	}

	files, err := ioutil.ReadDir(filepath.Join(`.`, path, `templates`))
	if err != nil {
		log.Fatal(err)
	}

	templates := []template.Template{}

	for _, f := range files {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)

		if languages[template.Canon(base)] {
			templates = append(templates, template.Template{Name: base, Type_: ext})
		}
	}

	orders, err := order.ReadOrder(`./config/gi/cache/templates/order`)
	if err != nil {
		log.Fatal(err)
	}

	o := template.Orderer{
		Templates: templates,
		Special:   orders,
	}
	o = template.Sort(o)

	if err := readFile(os.Stdout, o.Templates...); err != nil {
		log.Fatal(err)
	}
}

func readFile(w io.Writer, files ...template.Template) error {
	for _, file := range files {
		err := func(name, ext string) error {
			file, err := os.Open(filepath.Join(`.`, path, `templates`, name+ext))
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
		}(file.Name, file.Type_)

		if err != nil {
			return err
		}
	}

	return nil
}
