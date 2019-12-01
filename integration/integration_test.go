// +build integration

package integration

import (
	"bufio"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestMain(m *testing.M) {
	if err := exec.Command("make", "-C", "../", "build").Run(); err != nil {
		log.Fatalf("Fail to run make: %v", err)
	}

	code := m.Run()

	if err := exec.Command("make", "-C", "../", "clean").Run(); err != nil {
		log.Fatalf("Fail to run make: %v", err)
	}

	os.Exit(code)
}

func TestCli(t *testing.T) {
	actual, err := exec.Command("../gi", "go").CombinedOutput()
	assert.NoError(t, err)

	goldenPath := `./testdata/cli.golden`

	if *update {
		require.NoError(t, ioutil.WriteFile(goldenPath, actual, 0644))
	}

	expected, err := ioutil.ReadFile(goldenPath)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestCheckGitIgnoreIO(t *testing.T) {
	resp, err := http.Get(`https://www.gitignore.io/api/androidstudio,java,go,ada,zsh,c,gradle`)
	require.NoError(t, err)
	defer resp.Body.Close()

	expected := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan(); i++ {
		if i < 4 {
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

	actual, err := exec.Command("../gi", "androidstudio", "java", "go", "ada", "zsh", "c", "gradle", "go").CombinedOutput()
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, actual)
}
