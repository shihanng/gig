package order_test

import (
	"testing"

	"github.com/shihanng/gig/internal/order"
	"github.com/stretchr/testify/assert"
)

func TestReadOrder(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name      string
		args      args
		want      map[string]int
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "happy case",
			args: args{
				path: `./testdata/order`,
			},
			want: map[string]int{
				"java":          0,
				"gradle":        1,
				"androidstudio": 2,
				"visualstudio":  3,
				"umbraco":       4,
			},
			assertion: assert.NoError,
		},
		{
			name: "not found",
			args: args{
				path: `./testdata/unknown`,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.ReadOrder(tt.args.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
