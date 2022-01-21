//go:build integration

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/shihanng/gig/cmd"
	"github.com/stretchr/testify/suite"
)

var update = flag.Bool("update", false, "update .golden files")

type MainTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *MainTestSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "gig")
	s.Require().NoError(err)
	s.tempDir = dir
}

func (s *MainTestSuite) TearDownTest() {
	s.Require().NoError(os.RemoveAll(s.tempDir))
}

func (s *MainTestSuite) TestGen() {
	os.Args = []string{"gig", "gen", "go", "--cache-path", s.tempDir}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	goldenPath := `./testdata/gen.golden`

	if *update {
		s.Require().NoError(ioutil.WriteFile(goldenPath, actual.Bytes(), 0644))
	}

	expected, err := ioutil.ReadFile(goldenPath)
	s.Require().NoError(err)
	s.Assert().Equal(expected, actual.Bytes())
}

func (s *MainTestSuite) TestCheckGitIgnoreIO() {
	// Testing against "reactnative", "mean" is avoided because the result for stack from
	// gitignore.io seems not in order.
	resp, err := http.Get(`https://www.toptal.com/developers/gitignore/api/django,androidstudio,java,go,ada,zsh,c,gradle`)
	s.Require().NoError(err)

	defer resp.Body.Close()

	expected := new(bytes.Buffer)
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan(); i++ {
		if i < 3 {
			continue
		}

		content := scanner.Text()

		if strings.HasPrefix(content, `# End of https://www.toptal.com/developers/gitignore/api/django,androidstudio,java,go,ada,zsh,c,gradle`) {
			break
		}

		_, err := expected.WriteString(content + "\n")
		s.Require().NoError(err)
	}

	expectedBytes := expected.Bytes()
	expectedBytes = expectedBytes[:len(expectedBytes)-1]

	os.Args = []string{"gig", "--cache-path", s.tempDir, "gen",
		"Django", "androidstudio", "java", "go", "ada", "zsh", "c", "gradle", "go"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	s.Assert().Equal(expectedBytes, actual.Bytes())
}

func (s *MainTestSuite) TestList() {
	resp, err := http.Get(`https://www.toptal.com/developers/gitignore/api/list`)
	s.Require().NoError(err)

	defer resp.Body.Close()

	expected, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)

	expectedS := strings.Split(strings.ReplaceAll(string(expected), ",", "\n"), "\n")

	os.Args = []string{"gig", "--cache-path", s.tempDir, "list"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	actualS := strings.Split(strings.ToLower(actual.String()), "\n")
	actualS = actualS[:len(actualS)-1]

	s.Assert().Equal(expectedS, actualS)
}

func (s *MainTestSuite) TestVersion() {
	os.Args = []string{"gig", "--cache-path", s.tempDir, "version", "-c", "f0bddaeda3368130d52bde2b62a9df741f6117d4"}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	expected := strings.Join([]string{
		"gig version test",
		fmt.Sprintf("Cached github.com/toptal/gitignore in: %s", s.tempDir),
		"Using github.com/toptal/gitignore commit hash: f0bddaeda3368130d52bde2b62a9df741f6117d4",
	}, "\n") + "\n"

	s.Assert().Equal(expected, actual.String())
}

func (s *MainTestSuite) TestAutogen() {
	os.Args = []string{"gig", "autogen", "--cache-path", s.tempDir}

	actual := new(bytes.Buffer)

	cmd.Execute(actual, "test")

	goldenPath := `./testdata/autogen.golden`

	if *update {
		s.Require().NoError(ioutil.WriteFile(goldenPath, actual.Bytes(), 0644))
	}

	expected, err := ioutil.ReadFile(goldenPath)
	s.Require().NoError(err)
	s.Assert().Equal(expected, actual.Bytes())
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
