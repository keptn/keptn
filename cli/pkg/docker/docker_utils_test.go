package docker

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
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

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

var imageAvailabilityTests = []struct {
	image string
	tag   string
	err   error
}{
	{"docker.io/keptn/installer", "0.6.1", nil},
	{"docker.io/keptn/installer", "-1", errors.New("Provided image not found: Tag not found")},
	{"quay.io/keptn/installer", "1", errors.New("Provided image not found: 401 Unauthorized")},
	{"keptn/installer", "0.6.1", nil},
}

func TestCheckImageAvailablity(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		if strings.Contains(req.URL.String(), "docker.io") {
			if strings.Contains(req.URL.String(), "/tags/0.6.1") {
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			}
			return &http.Response{
				StatusCode: 404,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`Tag not found`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}

		} else if strings.Contains(req.URL.String(), "quay.io") {
			return &http.Response{
				StatusCode: 401,
				Status:     "401 Unauthorized",
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(``)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	for _, tt := range imageAvailabilityTests {
		t.Run(tt.image, func(t *testing.T) {
			err := CheckImageAvailability(tt.image, tt.tag, client)
			if err != tt.err {
				if !(err != nil && tt.err != nil && err.Error() == tt.err.Error()) {
					t.Errorf("got %q, want %q", err, tt.err)
				}
			}
		})
	}
}
