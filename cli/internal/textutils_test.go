package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JSONPathToJSONObj(t *testing.T) {

	var tests = []struct {
		name string
		in   []string
		out  string
		err  bool
	}{
		{"0", []string{"a.b=v1", "a.b.c=v2"}, "", true},
		{"1", []string{""}, "", true},
		{"2", []string{"a.b.c=v1"}, `{"a":{"b":{"c":"v1"}}}`, false},
		{"3", []string{"a.b.c=v", "a.b.d=v2"}, `{"a":{"b":{"c":"v","d":"v2"}}}`, false},
		{"4", []string{"a.b.c=v", "a.b.d=v2", "b.c.d.e=v3"}, `{"a":{"b":{"c":"v","d":"v2"}},"b":{"c":{"d":{"e":"v3"}}}}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := JSONPathToJSONObj(tt.in)
			assert.Equal(t, tt.out, out)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}

func TestUnfoldMap(t *testing.T) {
	type args struct {
		toAdd map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "minimal event data",
			args: args{
				toAdd: map[string]string{"project": "pr", "stage": "st", "service": "sv"},
			},
			want:    map[string]interface{}{"project": "pr", "stage": "st", "service": "sv"},
			wantErr: assert.NoError,
		},
		{
			name: "valid additional fields",
			args: args{
				toAdd: map[string]string{"project": "pr", "stage": "st", "service": "sv", "a.b": "value", "b.a": "othervalue"},
			},
			want:    map[string]interface{}{"project": "pr", "stage": "st", "service": "sv", "a": map[string]interface{}{"b": "value"}, "b": map[string]interface{}{"a": "othervalue"}},
			wantErr: assert.NoError,
		},
		{
			name: "valid additional fields",
			args: args{
				toAdd: map[string]string{"project": "pr", "stage": "st", "service": "sv", "a.b": "value", "b.a": "othervalue", "b.a.c": "myotherCvalue"},
			},
			want:    map[string]interface{}{},
			wantErr: assert.Error,
		},
		{
			name: "receive nil argument",
			args: args{
				toAdd: nil,
			},
			want:    map[string]interface{}{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnfoldMap(tt.args.toAdd)
			if !tt.wantErr(t, err, fmt.Sprintf("AddToData(%v)", tt.args.toAdd)) {
				return
			}
			assert.Equalf(t, tt.want, got, "AddToData(%v)", tt.args.toAdd)
		})
	}
}
