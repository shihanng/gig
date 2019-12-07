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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/gi/internal/file"
	"github.com/shihanng/gi/internal/order"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

const sourceRepo = `https://github.com/toptal/gitignore.git`

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen [template name]",
	Short: "Generates .gitignore of the given inputs",
	Long: `At the very first run the program will clone the templates
repository https://github.com/toptal/gitignore.git into $XDG_CACHE_HOME/gi.
This means that internet connection is not required after the first successful run.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := git.PlainClone(templatePath, false, &git.CloneOptions{
			URL:      sourceRepo,
			Depth:    1,
			Progress: ioutil.Discard,
		})
		if err != nil && err != git.ErrRepositoryAlreadyExists {
			return errors.Wrap(err, "gen: git clone")
		}

		languages := make(map[string]bool, len(args))

		for _, arg := range args {
			languages[file.Canon(arg)] = true
		}

		files, err := file.Filter(filepath.Join(templatePath, `templates`), languages)
		if err != nil {
			return err
		}

		orders, err := order.ReadOrder(filepath.Join(templatePath, `templates`, `order`))
		if err != nil {
			return err
		}

		files = file.SortFiles(files, orders)

		if err := file.Compose(os.Stdout, filepath.Join(templatePath, `templates`), files...); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
