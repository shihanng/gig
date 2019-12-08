package repo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const testSourceRepo = `https://github.com/thockin/test.git`

type RepoSuite struct {
	suite.Suite
	tempDir string
}

func (s *RepoSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "gi")
	require.NoError(s.T(), err)
	s.tempDir = dir
}

func (s *RepoSuite) TearDownTest() {
	require.NoError(s.T(), os.RemoveAll(s.tempDir))
}

// Order of the test cases is important
func (s *RepoSuite) TestNew() {
	type args struct {
		path       string
		repoSource string
	}

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "invalid source",
			args: args{
				path:       s.tempDir,
				repoSource: "",
			},
			assertion: assert.Error,
		},
		{
			name: "new clone",
			args: args{
				path:       s.tempDir,
				repoSource: testSourceRepo,
			},
			assertion: assert.NoError,
		},
		{
			name: "repo exists",
			args: args{
				path:       s.tempDir,
				repoSource: testSourceRepo,
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.path, tt.args.repoSource)
			tt.assertion(t, err)
		})
	}
}

func (s *RepoSuite) TestCheckout() {
	repo, err := New(s.tempDir, testSourceRepo)
	s.Require().NoError(err)

	type args struct {
		commitHash string
	}

	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "master",
			args: args{
				commitHash: "",
			},
			want:      "0d23e960376c28f21bbae5966a603e2f90ca2785",
			assertion: assert.NoError,
		},
		{
			name: "first commit",
			args: args{
				commitHash: "07ec2094347ba9cd825dce20909c215ca0dc6f37",
			},
			want:      "07ec2094347ba9cd825dce20909c215ca0dc6f37",
			assertion: assert.NoError,
		},
		{
			name: "undefined commit",
			args: args{
				commitHash: "58e32169bcb1b615cc8f4820e0299d07c6a679d2",
			},
			want:      "",
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := Checkout(repo, tt.args.commitHash)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
