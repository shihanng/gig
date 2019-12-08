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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

var rootCmd = &cobra.Command{
	Use:   "gi",
	Short: "A tool that generates .gitignore",
	Long: `gi is a command line tool to help you create useful .gitignore files
for your project. It is inspired by gitignore.io and make use of
the large collection of useful .gitignore templates of the web service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := git.PlainClone(templatePath, false, &git.CloneOptions{
			URL:      sourceRepo,
			Depth:    1,
			Progress: ioutil.Discard,
		})
		if err != nil && err != git.ErrRepositoryAlreadyExists {
			return errors.Wrap(err, "root: git clone")
		}

		return nil
	},
}

var templatePath string

func init() {
	cobra.OnInitialize(func() {
		templatePath = filepath.Join(xdg.CacheHome(), `gi`)
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
