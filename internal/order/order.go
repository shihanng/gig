// Package order provides helper function to parse
// the order file in https://github.com/toptal/gitignore
package order

import (
	"bufio"
	"os"

	"github.com/cockroachdb/errors"
)

// ReadOrder parses the order find in the provided path and
// returns the order of each items in the file. For the following content file
//
//	# A comment
//	go
//
//	elm
//
// We should get the following
//
//	"go": 0
//	"elm": 1
func ReadOrder(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "order: open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	orders := make(map[string]int)

	for n := 0; scanner.Scan(); {
		line := scanner.Text()
		if line != "" && !isComment(line) {
			orders[line] = n
			n++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "order: scanning")
	}

	return orders, nil
}

func isComment(line string) bool {
	return line != "" && line[0] == '#'
}
