package common

import (
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"time"
)

type RollbackFunc func() error

func ParseTimestamp(ts string, theClock clock.Clock) time.Time {
	parsedTime, err := timeutils.ParseTimestamp(ts)
	if err != nil {
		if theClock == nil {
			return time.Now().UTC()
		} else {
			return theClock.Now().UTC()
		}
	}
	return *parsedTime
}

func Stringp(s string) *string {
	return &s
}

func Merge(in1, in2 interface{}) interface{} {
	switch in1 := in1.(type) {
	case []interface{}:
		in2, ok := in2.([]interface{})
		if !ok {
			return in1
		}
		return append(in1, in2...)
	case map[string]interface{}:
		in2, ok := in2.(map[string]interface{})
		if !ok {
			return in1
		}
		for k, v2 := range in2 {
			if v1, ok := in1[k]; ok {
				in1[k] = Merge(v1, v2)
			} else {
				in1[k] = v2
			}
		}
	case nil:
		in2, ok := in2.(map[string]interface{})
		if ok {
			return in2
		}
	}
	return in1
}
