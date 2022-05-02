package v1

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromSequenceExecution(t *testing.T) {
	type args struct {
		se models.SequenceExecution
	}
	tests := []struct {
		name string
		args args
		want JsonStringEncodedSequenceExecution
	}{
		{
			name: "transform sequence execution",
			args: args{
				se: testSequenceExecution,
			},
			want: testJsonStringEncodedSequenceExecution,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := ModelTransformer{}
			got := mt.TransformToDBModel(tt.args.se)
			require.Equal(t, tt.want, got)
		})
	}
}
