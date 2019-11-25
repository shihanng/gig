package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
