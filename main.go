package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/shihanng/gi/internal/file"
	"github.com/shihanng/gi/internal/order"
	"gopkg.in/src-d/go-git.v4"
)

const sourceRepo = `https://github.com/toptal/gitignore.git`

func main() {
	path := filepath.Join(xdg.CacheHome(), `gi`)

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

	files, err := file.Filter(filepath.Join(path, `templates`), languages)
	if err != nil {
		log.Fatal(err)
	}

	orders, err := order.ReadOrder(filepath.Join(path, `templates`, `order`))
	if err != nil {
		log.Fatal(err)
	}

	files = file.Sort(files, orders)

	if err := file.Compose(os.Stdout, filepath.Join(path, `templates`), files...); err != nil {
		log.Fatal(err)
	}
}
