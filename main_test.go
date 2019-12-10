// +build integration

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/shihanng/gig/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestCli(t *testing.T) {
	os.Args = []string{"gig", "gen", "go"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	goldenPath := `./testdata/cli.golden`

	if *update {
		require.NoError(t, ioutil.WriteFile(goldenPath, actual.Bytes(), 0644))
	}

	expected, err := ioutil.ReadFile(goldenPath)
	require.NoError(t, err)
	assert.Equal(t, expected, actual.Bytes())
}

func TestCheckGitIgnoreIO(t *testing.T) {
	// Testing against "reactnative", "mean" is avoided because the result for stack from
	// gitignore.io seems not in order.
	resp, err := http.Get(`https://www.gitignore.io/api/django,androidstudio,java,go,ada,zsh,c,gradle`)
	require.NoError(t, err)

	defer resp.Body.Close()

	expected := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan(); i++ {
		if i < 3 {
			continue
		}

		content := scanner.Text()

		if strings.HasPrefix(content, `# End of https://www.gitignore.io/api/`) {
			break
		}

		_, err := expected.WriteString(content + "\n")
		require.NoError(t, err)
	}

	expectedBytes := expected.Bytes()
	expectedBytes = expectedBytes[:len(expectedBytes)-1]

	os.Args = []string{"gig", "gen", "Django", "androidstudio", "java", "go", "ada", "zsh", "c", "gradle", "go"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	assert.Equal(t, string(expectedBytes), actual.String())
}

func TestList(t *testing.T) {
	resp, err := http.Get(`https://www.gitignore.io/api/list`)
	require.NoError(t, err)

	defer resp.Body.Close()

	expected, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedS := bytes.Split(bytes.ReplaceAll(expected, []byte(","), []byte("\n")), []byte("\n"))

	os.Args = []string{"gig", "-c", "640f03b1f9906c5dcb788d36ec5c1095264a10ae", "list"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	actualS := bytes.Split(bytes.ToLower(actual.Bytes()), []byte("\n"))
	actualS = actualS[:len(actualS)-1]

	assert.Equal(t, expectedS, actualS)
}

func TestVersion(t *testing.T) {
	os.Args = []string{"gig", "version", "-c", "f0bddaeda3368130d52bde2b62a9df741f6117d4"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	expected := strings.Join([]string{
		"gig version test",
		fmt.Sprintf("Cached github.com/toptal/gitignore in: %s", filepath.Join(xdg.CacheHome(), `gig`)),
		"Using github.com/toptal/gitignore commit hash: f0bddaeda3368130d52bde2b62a9df741f6117d4",
	}, "\n") + "\n"

	assert.Equal(t, expected, actual.String())
}
