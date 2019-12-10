// +build integration

package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/shihanng/gig/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	cmd.Execute(actual)

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

	cmd.Execute(actual)

	actualS := bytes.Split(bytes.ToLower(actual.Bytes()), []byte("\n"))
	actualS = actualS[:len(actualS)-1]

	assert.Equal(t, expectedS, actualS)
}
