package eventmatcher

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEventMatcherUnableToDecodeEventData(t *testing.T) {
	require.False(t, EventMatcher{}.Matches(models.KeptnContextExtendedCE{Data: 0}))
}

func TestEventMatcher_Matches(t *testing.T) {
	type fields struct {
		Project string
		Stage   string
		Service string
	}
	type args struct {
		e models.KeptnContextExtendedCE
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no filter",
			fields: fields{
				Project: "",
				Stage:   "",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "exact filter",
			fields: fields{
				Project: "pr1",
				Stage:   "st1",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter project",
			fields: fields{
				Project: "pr1",
				Stage:   "",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter stage",
			fields: fields{
				Project: "",
				Stage:   "st1",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter service",
			fields: fields{
				Project: "",
				Stage:   "",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter project stage",
			fields: fields{
				Project: "pr1",
				Stage:   "st1",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter project service",
			fields: fields{
				Project: "pr1",
				Stage:   "",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "partial filter stage service",
			fields: fields{
				Project: "",
				Stage:   "st1",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: true,
		},
		{
			name: "full filter project - mismatch",
			fields: fields{
				Project: "pr2",
				Stage:   "st1",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
		{
			name: "full filter stage - mismatch",
			fields: fields{
				Project: "pr1",
				Stage:   "st2",
				Service: "sv1",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
		{
			name: "full filter service - mismatch",
			fields: fields{
				Project: "pr1",
				Stage:   "st1",
				Service: "sv2",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
		{
			name: "partial filter project - mismatch",
			fields: fields{
				Project: "pr2",
				Stage:   "",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
		{
			name: "partial filter stage - mismatch",
			fields: fields{
				Project: "",
				Stage:   "st2",
				Service: "",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
		{
			name: "partial filter service - mismatch",
			fields: fields{
				Project: "",
				Stage:   "",
				Service: "sv2",
			},
			args: args{
				e: models.KeptnContextExtendedCE{Data: v0_2_0.EventData{
					Project: "pr1",
					Stage:   "st1",
					Service: "sv1",
				}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := EventMatcher{
				Project: tt.fields.Project,
				Stage:   tt.fields.Stage,
				Service: tt.fields.Service,
			}
			if got := ef.Matches(tt.args.e); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
