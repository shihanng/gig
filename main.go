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

	"github.com/shihanng/gi/internal/file"
	"github.com/shihanng/gi/internal/order"
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
		languages[file.Canon(arg)] = true
	}

	files, err := ioutil.ReadDir(filepath.Join(`.`, path, `templates`))
	if err != nil {
		log.Fatal(err)
	}

	giFiles := []file.File{}

	for _, f := range files {
		filename := f.Name()
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)

		if languages[file.Canon(base)] {
			giFiles = append(giFiles, file.File{Name: base, Typ: ext})
		}
	}

	orders, err := order.ReadOrder(`./config/gi/cache/templates/order`)
	if err != nil {
		log.Fatal(err)
	}

	giFiles = file.Sort(giFiles, orders)

	if err := readFile(os.Stdout, giFiles...); err != nil {
		log.Fatal(err)
	}
}

func readFile(w io.Writer, files ...file.File) error {
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
		}(file.Name, file.Typ)

		if err != nil {
			return err
		}
	}

	return nil
}
