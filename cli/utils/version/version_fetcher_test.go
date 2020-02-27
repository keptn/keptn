package version

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

// testingHTTPClient builds a test client with a httptest server
func testingHTTPClient(handler http.Handler) (*http.Client, string, func()) {
	server := httptest.NewTLSServer(handler)

	cert, err := x509.ParseCertificate(server.TLS.Certificates[0].Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	certpool := x509.NewCertPool()
	certpool.AddCert(cert)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				RootCAs: certpool,
			},
		},
	}

	return client, server.URL, server.Close
}

func TestGetCLIVersionInfo(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expect GET request")
		assert.Equal(t, r.Header.Get("user-agent"), "keptn/cli:0.6.0", "Expect user-agent header")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{ "cli": { "stable": ["0.5.2", "0.6.0"], "prerelease": ["0.6.0-beta2"] } }`)
	})

	httpClient, url, teardown := testingHTTPClient(handler)
	defer teardown()

	client := newVersionFetcherClient()
	client.httpClient = httpClient
	client.versionUrl = url

	cliVersionInfo, err := client.getCLIVersionInfo("0.6.0")
	assert.Equal(t, err, nil, "Received unexpected error")
	assert.Equal(t, cliVersionInfo.Prerelease, []string{"0.6.0-beta2"}, "Received unexpected content")
	assert.Equal(t, cliVersionInfo.Stable, []string{"0.5.2", "0.6.0"}, "Received unexpected content")
}
