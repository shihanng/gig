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
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/gig/internal/file"
	"github.com/spf13/cobra"
)

func newSearchCmd(c *command) *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search and select supported templates list",
		Long: `Search and select supported templates list.

This subcommand depends on fzf (https://github.com/junegunn/fzf)
for the search functionality.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			templates, err := file.List(c.templatePath())
			if err != nil {
				return err
			}
			templateStr := strings.Join(templates, "\n")

			var selected bytes.Buffer

			cmdName, cmdArgs := "sh", []string{"-c"}
			if runtime.GOOS == "windows" {
				cmdName = "cmd"
				cmdArgs[0] = "/c"
			}

			shCmd := exec.Command(cmdName, append(cmdArgs, c.searchTool)...)

			shCmd.Stdin = strings.NewReader(templateStr)
			shCmd.Stdout = &selected
			shCmd.Stderr = os.Stderr

			if err := shCmd.Run(); err != nil {
				var ee *exec.ExitError
				if errors.As(err, &ee) && ee.ExitCode() == 130 {
					return nil
				}

				return errors.Wrap(err, "cmd: searching")
			}

			scanner := bufio.NewScanner(&selected)

			var items []string
			for scanner.Scan() {
				items = append(items, strings.TrimSpace(scanner.Text()))
			}

			if err := scanner.Err(); err != nil {
				return errors.Wrap(err, "cmd: read search result")
			}

			return c.generateIgnoreFile(items)
		},
	}
}
