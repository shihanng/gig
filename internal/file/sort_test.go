package file_test

import (
	"testing"

	"github.com/shihanng/gig/internal/file"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	type args struct {
		items        []string
		specialOrder map[string]int
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nothing",
		},
		{
			name: "only text",
			args: args{
				items: []string{"Z", "a", "f", "z", "A"},
			},
			want: []string{"a", "A", "f", "Z", "z"},
		},
		{
			name: "special",
			args: args{
				items: []string{
					"Ada",
					"AndroidStudio",
					"C",
					"Go",
					"Gradle",
					"Java",
					"Zsh",
				},
				specialOrder: map[string]int{
					"androidstudio": 2,
					"gradle":        1,
					"java":          0,
					"umbraco":       4,
					"visualstudio":  3,
				},
			},
			want: []string{
				"Ada",
				"C",
				"Go",
				"Java",
				"Zsh",
				"Gradle",
				"AndroidStudio",
			},
		},
		{
			name: "special 2",
			args: args{
				items: []string{"Django", "androidstudio", "java", "go", "ada", "zsh", "c", "gradle", "go"},
				specialOrder: map[string]int{
					"androidstudio": 2,
					"gradle":        1,
					"java":          0,
					"umbraco":       4,
					"visualstudio":  3,
				},
			},
			want: []string{"ada", "c", "Django", "go", "go", "java", "zsh", "gradle", "androidstudio"},
		},
		{
			name: "special 3",
			args: args{
				items: []string{"Django", "androidstudio", "ada", "zsh", "c", "gradle", "go"},
				specialOrder: map[string]int{
					"androidstudio": 2,
					"gradle":        1,
					"java":          0,
					"umbraco":       4,
					"visualstudio":  3,
				},
			},
			want: []string{"ada", "c", "Django", "go", "zsh", "gradle", "androidstudio"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, file.Sort(tt.args.items, tt.args.specialOrder))
		})
	}
}
