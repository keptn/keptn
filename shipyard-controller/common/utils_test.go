package common

import (
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

func TestChangeEventType(t *testing.T) {
	type args struct {
		current   string
		wanted    string
		delimiter string
	}

	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				current:   "some.funny.event.triggered",
				wanted:    "waiting",
				delimiter: ".",
			},
			want: "some.funny.event.waiting",
		},
		{
			args: args{
				current:   "some,funny,event,triggered",
				wanted:    "waiting",
				delimiter: ",",
			},
			want: "some,funny,event,waiting",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := ChangeEventType(tt.args.current, tt.args.wanted, tt.args.delimiter); got != tt.want {
				t.Errorf("ChangeEventType() = %v, want %v", got, tt.want)
			}
		})
	}

}
