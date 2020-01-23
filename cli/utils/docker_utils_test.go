package utils

import "testing"

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
