package docker

import (
	"errors"
	"testing"
)

var imageSplitTests = []struct {
	in       string
	outImage string
	outTag   string
}{
	{"docker.io/keptn/installer:0.6.0.beta", "docker.io/keptn/installer", "0.6.0.beta"},
	{"docker.io/keptn/installer", "docker.io/keptn/installer", "latest"},
	{"keptn/installer:0.6.0.beta", "keptn/installer", "0.6.0.beta"},
	{"keptn/installer", "keptn/installer", "latest"},
	{"installer:0.6.0.beta", "installer", "0.6.0.beta"},
	{"installer", "installer", "latest"},
}

func TestSplitImageName(t *testing.T) {
	for _, tt := range imageSplitTests {
		t.Run(tt.in, func(t *testing.T) {
			outImage, outTag := SplitImageName(tt.in)
			if outImage != tt.outImage {
				t.Errorf("got %q, want %q", outImage, tt.outImage)
			}
			if outTag != tt.outTag {
				t.Errorf("got %q, want %q", outTag, tt.outTag)
			}
		})
	}
}

var imageAvailabilityTests = []struct {
	image string
	tag   string
	err   error
}{
	{"docker.io/keptn/installer", "0.6.1", nil},
	{"docker.io/keptn/installer", "-1", errors.New("Provided image not found: Tag not found")},
	{"quay.io/keptn/installer", "-1", errors.New("Provided image not found: 401 Unauthorized")},
	{"keptn/installer", "0.6.1", nil},
}

func TestCheckImageAvailablity(t *testing.T) {
	for _, tt := range imageAvailabilityTests {
		t.Run(tt.image, func(t *testing.T) {
			err := CheckImageAvailability(tt.image, tt.tag)
			if err != tt.err {
				if !(err != nil && tt.err != nil && err.Error() == tt.err.Error()) {
					t.Errorf("got %q, want %q", err, tt.err)
				}
			}
		})
	}
}
