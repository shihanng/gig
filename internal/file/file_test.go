package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
