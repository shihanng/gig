package integration

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
