package lib_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
)

func TestCurlValidator_ResolveIPAddresses(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		ipResolver lib.IPResolver
		want       []string
	}{
		{
			name: "error output",
			url:  "http://some-url",
			ipResolver: lib.IPResolver{
				LookupIP: func(host string) ([]net.IP, error) {
					return make([]net.IP, 0), fmt.Errorf("some error")
				},
			},
			want: make([]string, 0),
		},
		{
			name: "no existing address",
			url:  "http://some-url",
			ipResolver: lib.IPResolver{
				LookupIP: func(host string) ([]net.IP, error) {
					return make([]net.IP, 0), nil
				},
			},
			want: make([]string, 0),
		},
		{
			name: "ip addresses list",
			url:  "http://some-url",
			ipResolver: lib.IPResolver{
				LookupIP: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")}, nil
				},
			},
			want: []string{"1.1.1.1", "2.2.2.2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ipResolver.ResolveIPAdresses(tt.url)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDeniedURLS(t *testing.T) {
	kubeEnvs := map[string]string{"KUBERNETES_SERVICE_HOST": "1.2.3.4", "KUBERNETES_SERVICE_PORT": "9876"}
	urls := lib.GetDeniedURLs(kubeEnvs)

	expected := []string{"1.2.3.4", "kubernetes:9876", "kubernetes.default:9876", "kubernetes.default.svc:9876", "kubernetes.default.svc.cluster.local:9876", "1.2.3.4:9876"}

	require.Equal(t, 6, len(urls))
	require.Equal(t, expected, urls)
}
