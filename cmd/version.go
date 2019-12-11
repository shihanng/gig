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

	"github.com/spf13/cobra"
)

func newVersionCmd(c *command) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number and other useful info",
		Run:   c.versionRunE,
	}
}

func (c *command) versionRunE(cmd *cobra.Command, args []string) {
	fmt.Fprintf(c.output, "gig version %s\n", c.version)
	fmt.Fprintf(c.output, "Cached github.com/toptal/gitignore in: %s\n", c.templatePath)
	fmt.Fprintf(c.output, "Using github.com/toptal/gitignore commit hash: %s\n", c.commitHash)
}
