package common

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
)

func TestParseTimestamp(t *testing.T) {
	correctISO8601Timestamp := "2020-01-02T15:04:05.000Z"

	timeObj, _ := time.Parse(timeutils.KeptnTimeFormatISO8601, correctISO8601Timestamp)

	mockClock := clock.NewMock()

	type args struct {
		ts       string
		theClock clock.Clock
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "compatible timestamp provided",
			args: args{
				ts:       correctISO8601Timestamp,
				theClock: nil,
			},
			want: timeObj,
		},
		{
			name: "incompatible timestamp provided - return now",
			args: args{
				ts:       "invalid",
				theClock: mockClock,
			},
			want: mockClock.Now().UTC(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTimestamp(tt.args.ts, tt.args.theClock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	type args struct {
		in1 interface{}
		in2 interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "merge maps",
			args: args{
				in1: map[string]interface{}{
					"foo": "bar",
				},
				in2: map[string]interface{}{
					"bar": "foo",
				},
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": "foo",
			},
		},
		{
			name: "merge map - string property of in2 should override property with same name in in1",
			args: args{
				in1: map[string]interface{}{
					"foo": "bar",
				},
				in2: map[string]interface{}{
					"foo": "foo",
				},
			},
			want: map[string]interface{}{
				"foo": "foo",
			},
		},
		{
			name: "merge different structures",
			args: args{
				in1: map[string]interface{}{
					"foo": "bar",
				},
				in2: map[string]interface{}{
					"bar": map[string]interface{}{
						"bar": "foo",
					},
				},
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"bar": "foo",
				},
			},
		},
		{
			name: "merge different structures 2",
			args: args{
				in1: []interface{}{"item1", "item2"},
				in2: []interface{}{"item3"},
			},
			want: []interface{}{"item3", "item1", "item2"},
		},
		{
			name: "merge structures with different types for same property names: map vs string",
			args: args{
				in1: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "xyz",
					},
				},
				in2: map[string]interface{}{
					"foo": "bar",
				},
			},
			want: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "merge structures with different types for same property names: slice vs string",
			args: args{
				in1: map[string]interface{}{
					"foo": []interface{}{"bar"},
				},
				in2: map[string]interface{}{
					"foo": "bar",
				},
			},
			want: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "merge structures with different types for same property names: string vs slice",
			args: args{
				in1: map[string]interface{}{
					"foo": "bar",
				},
				in2: map[string]interface{}{
					"foo": []interface{}{"bar"},
				},
			},
			want: map[string]interface{}{
				"foo": []interface{}{"bar"},
			},
		},
		{
			name: "merge structures with different types: nil vs map",
			args: args{
				in1: nil,
				in2: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "xyz",
					},
				},
			},
			want: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "xyz",
				},
			},
		},
		{
			name: "merge structures with different types: nil vs string",
			args: args{
				in1: nil,
				in2: "foo",
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args.in1, tt.args.in2)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestCopyMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "get copied map",
			args: args{
				m: map[string]interface{}{
					"bar": map[string]interface{}{
						"bar": "foo",
					},
				},
			},
			want: map[string]interface{}{
				"bar": map[string]interface{}{
					"bar": "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
