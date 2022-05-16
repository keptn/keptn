package lib

import (
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurlValidator_ResolveIPAddresses(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		ipResolver ipResolver
		want       []string
	}{
		{
			name: "unparsable address",
			url:  "http://some-url",
			ipResolver: ipResolver{
				parse: func(rawURL string) (*url.URL, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			want: make([]string, 0),
		},
		{
			name: "lookupIP failed",
			url:  "http://some-url",
			ipResolver: ipResolver{
				parse: func(rawURL string) (*url.URL, error) {
					return &url.URL{
						Host: "some-url",
					}, nil
				},
				lookupIP: func(host string) ([]net.IP, error) {
					return make([]net.IP, 0), fmt.Errorf("some lookupIP error")
				},
			},
			want: make([]string, 0),
		},
		{
			name: "no existing address",
			url:  "http://some-url",
			ipResolver: ipResolver{
				parse: func(rawURL string) (*url.URL, error) {
					return &url.URL{
						Host: "some-url",
					}, nil
				},
				lookupIP: func(host string) ([]net.IP, error) {
					return make([]net.IP, 0), nil
				},
			},
			want: make([]string, 0),
		},
		{
			name: "ip addresses list",
			url:  "http://some-url",
			ipResolver: ipResolver{
				parse: func(rawURL string) (*url.URL, error) {
					return &url.URL{
						Host: "some-url",
					}, nil
				},
				lookupIP: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")}, nil
				},
			},
			want: []string{"1.1.1.1", "2.2.2.2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ipResolver.Resolve(tt.url)
			require.Equal(t, tt.want, got)
		})
	}
}
