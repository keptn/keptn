// +build !nokubectl

package cmd

import (
	"reflect"
	"testing"
)

var iskeptnVersions = []struct {
	platform             string
	expectedPlatformType reflect.Type
	expectedErr          error
}{
	{"OpenShift", reflect.TypeOf(newOpenShiftPlatform()), nil},
	{"kubernetes", reflect.TypeOf(newKubernetesPlatform()), nil},
}

func TestSetPlatform(t *testing.T) {
	for _, tt := range iskeptnVersions {
		t.Run(tt.platform, func(t *testing.T) {
			*installParams.PlatformIdentifier = tt.platform
			err := setPlatform()
			if err != tt.expectedErr {
				t.Errorf("got %t, want %t", err, tt.expectedErr)
			}
			if reflect.TypeOf(p) != tt.expectedPlatformType {
				t.Error("wrong type")
			}
		})
	}
}
