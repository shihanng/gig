/*
Copyright © 2019 Shi Han NG <shihanng@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/gig/internal/file"
	"github.com/spf13/cobra"
	"github.com/src-d/enry/v2"
)

func newAutogenCmd(c *command) *cobra.Command {
	return &cobra.Command{
		Use:   "autogen",
		Short: "(EXPERIMENTAL) I'm Feeling Lucky",
		Long: `(EXPERIMENTAL) Automatically generate .gitignore based on the content of the current directory.

This feature uses github.com/src-d/enry
(which is a port of GitHub's linguist library) to detect
the programming languages in the current working directory
and pass that to the gen command. Known limitation of this
feature is that it could not detect a framework.`,
		RunE: c.autogenRunE,
	}
}

// Heavily borrowed from:
// https://github.com/src-d/enry/blob/697929e1498cbdb7726a4d3bf4c48e706ee8c967/cmd/enry/main.go#L27
func (c *command) autogenRunE(cmd *cobra.Command, args []string) error {
	templates, err := file.List(c.templatePath())
	if err != nil {
		return err
	}

	supported := make(map[string]bool, len(templates))

	for _, t := range templates {
		supported[file.Canon(t)] = true
	}

	var found []string

	errWalk := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if !f.Mode().IsDir() && !f.Mode().IsRegular() {
			return nil
		}

		if enry.IsVendor(path) ||
			enry.IsDotFile(path) ||
			enry.IsDocumentation(path) ||
			enry.IsConfiguration(path) {
			if f.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if f.IsDir() {
			return nil
		}

		//nolint:gomnd
		content, err := readFile(path, 16*1024)
		if err != nil {
			return err
		}

		language := enry.GetLanguage(path, content)
		if language == enry.OtherLanguage {
			return nil
		}

		if supported[file.Canon(language)] {
			found = append(found, language)
		}

		return nil
	})

	if errWalk != nil {
		return errors.Wrap(err, "cmd: walking file")
	}

	return c.generateIgnoreFile(found)
}

func readFile(path string, limit int64) ([]byte, error) {
	if limit <= 0 {
		b, err := ioutil.ReadFile(path)

		return b, errors.Wrap(err, "cmd: read file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "cmd: open file")
	}

	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "cmd: get file stat")
	}

	size := st.Size()
	if limit > 0 && size > limit {
		size = limit
	}

	buf := bytes.NewBuffer(nil)
	buf.Grow(int(size))

	_, err = io.Copy(buf, io.LimitReader(f, limit))

	return buf.Bytes(), errors.Wrap(err, "cmd: copy to buffer")
}
