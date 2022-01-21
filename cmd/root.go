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
	"io"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"github.com/shihanng/gig/internal/file"
	"github.com/shihanng/gig/internal/order"
	"github.com/shihanng/gig/internal/repo"
	"github.com/spf13/cobra"
)

func Execute(w io.Writer, version string) {
	command := &command{
		output:     w,
		cachePath:  filepath.Join(xdg.CacheHome(), `gig`),
		version:    version,
		searchTool: "fzf -m",
	}

	rootCmd := newRootCmd(command)

	rootCmd.PersistentFlags().StringVarP(&command.commitHash, "commit-hash", "c", "",
		`use templates from a specific commit hash of
github.com/toptal/gitignore`)

	rootCmd.PersistentFlags().StringVarP(&command.cachePath, "cache-path", "", filepath.Join(xdg.CacheHome(), `gig`),
		`location where the content of github.com/toptal/gitignore
will be cached in`)

	genCmd := newGenCmd(command)

	genCmd.Flags().BoolVarP(&command.genIsFile, "file", "f", false,
		"if specified will create .gitignore file in the current working directory")

	searchCmd := newSearchCmd(command)

	searchCmd.Flags().BoolVarP(&command.genIsFile, "file", "f", false,
		"if specified will create .gitignore file in the current working directory")

	autogenCmd := newAutogenCmd(command)

	autogenCmd.Flags().BoolVarP(&command.genIsFile, "file", "f", false,
		"if specified will create .gitignore file in the current working directory")

	rootCmd.AddCommand(
		newListCmd(command),
		genCmd,
		newVersionCmd(command),
		searchCmd,
		autogenCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(c *command) *cobra.Command {
	return &cobra.Command{
		Use:   "gig",
		Short: "A tool that generates .gitignore",
		Long: `gig is a command line tool to help you create useful .gitignore files
for your project. It is inspired by gitignore.io and make use of
the large collection of useful .gitignore templates of the web service.`,
		PersistentPreRunE: c.rootRunE,
	}
}

type command struct {
	output     io.Writer
	commitHash string
	cachePath  string
	version    string
	searchTool string

	genIsFile bool
}

func (c *command) rootRunE(cmd *cobra.Command, args []string) error {
	r, err := repo.New(c.cachePath, repo.SourceRepo)
	if err != nil {
		return err
	}

	ch, err := repo.Checkout(r, c.commitHash)
	if err != nil {
		return err
	}

	c.commitHash = ch

	return nil
}

func (c *command) generateIgnoreFile(items []string) error {
	orders, err := order.ReadOrder(filepath.Join(c.templatePath(), `order`))
	if err != nil {
		return err
	}

	items = file.Sort(items, orders)

	wc, err := c.newWriteCloser()
	if err != nil {
		return err
	}

	defer wc.Close()

	return file.Generate(wc, c.templatePath(), items...)
}

func (c *command) newWriteCloser() (io.WriteCloser, error) {
	if c.genIsFile {
		f, err := os.Create(".gitignore")
		if err != nil {
			return nil, errors.Wrap(err, "cmd: create new file")
		}

		return f, nil
	}

	return ioutil.WriteNopCloser(c.output), nil
}

func (c *command) templatePath() string {
	return filepath.Join(c.cachePath, `templates`)
}
