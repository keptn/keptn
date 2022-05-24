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
		want       AdrDomainNameMapping
		wanterr    error
	}{
		{
			name: "unparsable address",
			url:  "http://some-url",
			ipResolver: ipResolver{
				parse: func(rawURL string) (*url.URL, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			want:    make(AdrDomainNameMapping, 0),
			wanterr: fmt.Errorf("some error"),
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
			want:    make(AdrDomainNameMapping, 0),
			wanterr: fmt.Errorf("some lookupIP error"),
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
			want: make(AdrDomainNameMapping, 0),
		},
		{
			name: "ip addresses list no host",
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
				lookupAddr: func(addr string) ([]string, error) {
					return []string{}, nil
				},
			},
			want: AdrDomainNameMapping{
				"1.1.1.1": {},
				"2.2.2.2": {},
			},
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
				lookupAddr: func(addr string) ([]string, error) {
					return []string{"myhost"}, nil
				},
			},
			want: AdrDomainNameMapping{
				"1.1.1.1": {"myhost"},
				"2.2.2.2": {"myhost"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ipResolver.Resolve(tt.url)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wanterr, err)
		})
	}
}
