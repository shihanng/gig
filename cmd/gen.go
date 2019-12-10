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
	"path/filepath"

	"github.com/shihanng/gig/internal/file"
	"github.com/shihanng/gig/internal/order"
	"github.com/spf13/cobra"
)

func newGenCmd(w io.Writer, templatePath string) *cobra.Command {
	return &cobra.Command{
		Use:   "gen [template name]",
		Short: "Generates .gitignore of the given inputs",
		Long: `Generates .gitignore of the given [template name]
which should contain one or more valid names (case insensitive).
Valid names can be obtained from the list subcommand.
At the very first run the program will clone the templates repository
https://github.com/toptal/gitignore.git into $XDG_CACHE_HOME/gig.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			items := args

			orders, err := order.ReadOrder(filepath.Join(templatePath, `templates`, `order`))
			if err != nil {
				return err
			}

			items = file.Sort(items, orders)

			return file.Generate(w, filepath.Join(templatePath, `templates`), items...)
		},
	}
}
