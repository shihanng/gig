package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/shihanng/gig/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestCli(t *testing.T) {
	os.Args = []string{"gig", "gen", "go"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual)

	goldenPath := `./testdata/cli.golden`

	if *update {
		require.NoError(t, ioutil.WriteFile(goldenPath, actual.Bytes(), 0644))
	}

	expected, err := ioutil.ReadFile(goldenPath)
	require.NoError(t, err)
	assert.Equal(t, expected, actual.Bytes())
}
