package file

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestFilter(t *testing.T) {
	type args struct {
		directory string
		filter    map[string]bool
	}

	tests := []struct {
		name      string
		args      args
		want      []File
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				directory: `testdata`,
				filter: map[string]bool{
					"go": true,
				},
			},
			want: []File{
				{Name: "Go", Typ: ".gitignore"},
			},
			assertion: assert.NoError,
		},
		{
			name: "not found",
			args: args{
				directory: `unknown`,
				filter: map[string]bool{
					"go": true,
				},
			},
			want:      []File{},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.args.directory, tt.args.filter)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompose(t *testing.T) {
	type args struct {
		files []File
	}

	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "single file",
			args: args{
				files: []File{{Name: "Go", Typ: ".gitignore"}},
			},
			wantW:     "Go.gitignore.golden",
			assertion: assert.NoError,
		},
		{
			name: "two files",
			args: args{
				files: []File{
					{Name: "Go", Typ: ".gitignore"},
					{Name: "C", Typ: ".gitignore"},
				},
			},
			wantW:     "GoC.gitignore.golden",
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, Compose(w, `testdata`, tt.args.files...))

			goldenPath := filepath.Join(`testdata`, tt.wantW)

			if *update {
				require.NoError(t, ioutil.WriteFile(goldenPath, w.Bytes(), 0644))
			}

			expected, err := ioutil.ReadFile(goldenPath)
			require.NoError(t, err)
			assert.Equal(t, expected, w.Bytes())
		})
	}
}

func TestSort(t *testing.T) {
	type args struct {
		f       []File
		special map[string]int
	}

	tests := []struct {
		name string
		args args
		want []File
	}{
		{
			name: "simple",
			args: args{
				f: []File{
					{Name: "B"},
					{Name: "A"},
				},
			},
			want: []File{
				{Name: "A"},
				{Name: "B"},
			},
		},
		{
			name: "special",
			args: args{
				f: []File{
					{Name: "B"},
					{Name: "A"},
				},
				special: map[string]int{
					"a": 2,
					"b": 1,
				},
			},
			want: []File{
				{Name: "B"},
				{Name: "A"},
			},
		},
		{
			name: "type",
			args: args{
				f: []File{
					{Name: "A", Typ: ".Patch"},
					{Name: "A", Typ: ".GitIgnore"},
				},
			},
			want: []File{
				{Name: "A", Typ: ".GitIgnore"},
				{Name: "A", Typ: ".Patch"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Sort(tt.args.f, tt.args.special))
		})
	}
}
