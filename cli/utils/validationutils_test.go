package utils

import (
	"testing"
)

var keptnVersions = []struct {
	in  string
	res bool
}{
	{"20191212.1033-latest", false},
	{"0.6.0.beta2", true},
	{"feature-443-20191213.1105", false},
	{"0.6.0.beta2-20191204.1329", false},
	{"0.6.0.beta2-201912044.1329", true},
}

func TestIsOfficialKeptnVersion(t *testing.T) {
	for _, tt := range keptnVersions {
		t.Run(tt.in, func(t *testing.T) {
			res := IsOfficialKeptnVersion(tt.in)
			if res != tt.res {
				t.Errorf("got %t, want %t", res, tt.res)
			}
		})
	}
}
