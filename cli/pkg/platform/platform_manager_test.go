// +build !nokubectl

package platform

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
			p, err := NewPlatformManager(tt.platform)
			if err != tt.expectedErr {
				t.Errorf("got %t, want %t", err, tt.expectedErr)
			}
			if reflect.TypeOf(p.platform) != tt.expectedPlatformType {
				t.Error("wrong type")
			}
		})
	}
}
