package file

import (
	"testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Sort(tt.args.items, tt.args.specialOrder))
		})
	}
}
