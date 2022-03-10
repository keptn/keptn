package common

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
)

type RollbackFunc func() error

// ParseTimestamp tries to parse the given timestamp.
// If for some reason, the provided value cannot be parsed, the current time is returned instead
// Optionally, the function allows to pass an implementation of the Clock interface, which is then used for determining the fallback value that should be returned in case
// the given timestamp could not be parsed
func ParseTimestamp(ts string, theClock clock.Clock) time.Time {
	parsedTime, err := timeutils.ParseTimestamp(ts)
	if err != nil {
		if theClock == nil {
			return time.Now().UTC()
		}
		return theClock.Now().UTC()
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
		mergeMaps(in1, in2)
	case string:
		if in2, ok := in2.(string); ok {
			return in2
		}
	case nil:
		in2, ok := in2.(map[string]interface{})
		if ok {
			return in2
		}
	}
	return in1
}

func mergeMaps(in1 map[string]interface{}, in2 map[string]interface{}) {
	for k, v2 := range in2 {
		if v1, ok := in1[k]; ok && v1 != nil {
			in1[k] = Merge(v1, v2)
		} else {
			in1[k] = v2
		}
	}
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}
