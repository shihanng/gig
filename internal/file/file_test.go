package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	type args struct {
		o Orderer
	}
	tests := []struct {
		name string
		args args
		want Orderer
	}{
		{
			name: "simple",
			args: args{
				o: Orderer{
					Files: []File{
						{Name: "B"},
						{Name: "A"},
					},
				},
			},
			want: Orderer{
				Files: []File{
					{Name: "A"},
					{Name: "B"},
				},
			},
		},
		{
			name: "special",
			args: args{
				o: Orderer{
					Files: []File{
						{Name: "B"},
						{Name: "A"},
					},
					Special: map[string]int{
						"a": 2,
						"b": 1,
					},
				},
			},
			want: Orderer{
				Files: []File{
					{Name: "B"},
					{Name: "A"},
				},
				Special: map[string]int{
					"a": 2,
					"b": 1,
				},
			},
		},
		{
			name: "type",
			args: args{
				o: Orderer{
					Files: []File{
						{Name: "A", Typ: ".Patch"},
						{Name: "A", Typ: ".GitIgnore"},
					},
				},
			},
			want: Orderer{
				Files: []File{
					{Name: "A", Typ: ".GitIgnore"},
					{Name: "A", Typ: ".Patch"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Sort(tt.args.o))
		})
	}
}
