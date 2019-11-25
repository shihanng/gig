package main

import (
	"log"
	"os"

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
	if err != nil {
		log.Fatal(err)
	}
}
