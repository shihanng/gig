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
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/shihanng/gi/internal/repo"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gi",
		Short: "A tool that generates .gitignore",
		Long: `gi is a command line tool to help you create useful .gitignore files
for your project. It is inspired by gitignore.io and make use of
the large collection of useful .gitignore templates of the web service.`,
	}
}

func rootCmdRun(templatePath string, commitHash *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		r, err := repo.New(templatePath, repo.SourceRepo)
		if err != nil {
			return err
		}

		_, err = repo.Checkout(r, *commitHash)
		if err != nil {
			return err
		}

		return nil
	}
}

func Execute() {
	templatePath := filepath.Join(xdg.CacheHome(), `gi`)

	var commitHash string

	rootCmd := newRootCmd()

	rootCmd.PersistentFlags().StringVarP(&commitHash, "commit-hash", "c", "",
		"use templates from a specific commit hash of github.com/toptal/gitignore")

	rootCmd.PersistentPreRunE = rootCmdRun(templatePath, &commitHash)

	rootCmd.AddCommand(
		newListCmd(templatePath),
		newGenCmd(templatePath),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
