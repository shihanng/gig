package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	type args struct {
		o orderer
	}
	tests := []struct {
		name string
		args args
		want orderer
	}{
		{
			name: "simple",
			args: args{
				o: orderer{
					templates: []template{
						{name: "B"},
						{name: "A"},
					},
				},
			},
			want: orderer{
				templates: []template{
					{name: "A"},
					{name: "B"},
				},
			},
		},
		{
			name: "special",
			args: args{
				o: orderer{
					templates: []template{
						{name: "B"},
						{name: "A"},
					},
					special: map[string]int{
						"a": 2,
						"b": 1,
					},
				},
			},
			want: orderer{
				templates: []template{
					{name: "B"},
					{name: "A"},
				},
				special: map[string]int{
					"a": 2,
					"b": 1,
				},
			},
		},
		{
			name: "type",
			args: args{
				o: orderer{
					templates: []template{
						{name: "A", type_: "Patch"},
						{name: "A", type_: "GitIgnore"},
					},
				},
			},
			want: orderer{
				templates: []template{
					{name: "A", type_: "GitIgnore"},
					{name: "A", type_: "Patch"},
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
