package main

import (
	"bufio"
	"fmt"
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
	readFile()
}

const filename = `Go.gitignore`

func readFile() {
	file, err := os.Open(filepath.Join(`.`, path, `templates`, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator
		fmt.Println(scanner.Text()) // token in unicode-char
	}
}
