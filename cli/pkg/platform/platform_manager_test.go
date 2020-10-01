// +build !nokubectl

package platform

import (
	"reflect"
	"testing"
)

var iskeptnVersions = []struct {
	isOpenShiftEnabled   bool
	expectedPlatformType reflect.Type
}{
	{true, reflect.TypeOf(newOpenShiftPlatform())},
	{false, reflect.TypeOf(newKubernetesPlatform())},
}

func TestSetPlatform(t *testing.T) {
	for _, tt := range iskeptnVersions {
		var testName string
		if tt.isOpenShiftEnabled {
			testName = "OpenShift"
		} else {
			testName = "Kubernetes"
		}
		t.Run(testName, func(t *testing.T) {
			p := NewPlatformManager(tt.isOpenShiftEnabled)
			if reflect.TypeOf(p.platform) != tt.expectedPlatformType {
				t.Error("wrong type")
			}
		})
	}
}
