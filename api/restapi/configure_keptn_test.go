package restapi

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getEnvConfig(t *testing.T) {
	defer os.Unsetenv("MAX_AUTH_ENABLED")
	defer os.Unsetenv("MAX_AUTH_REQUESTS_PER_SECOND")
	defer os.Unsetenv("MAX_AUTH_REQUESTS_BURST")
	defer os.Unsetenv("HIDE_DEPRECATED")
	defer os.Unsetenv("OAUTH_ENABLED")
	defer os.Unsetenv("OAUTH_PREFIX")
	_ = os.Setenv("MAX_AUTH_ENABLED", "false")
	_ = os.Setenv("MAX_AUTH_REQUESTS_PER_SECOND", "0.5")
	_ = os.Setenv("MAX_AUTH_REQUESTS_BURST", "1")
	_ = os.Setenv("HIDE_DEPRECATED", "true")
	_ = os.Setenv("OAUTH_ENABLED", "true")
	_ = os.Setenv("OAUTH_PREFIX", "prefix")

	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, false, config.MaxAuthEnabled)
	require.Equal(t, 0.5, config.MaxAuthRequestsPerSecond)
	require.Equal(t, 1, config.MaxAuthRequestBurst)
	require.Equal(t, true, config.HideDeprecated)
	require.Equal(t, true, config.OAuthEnabled)
	require.Equal(t, "prefix", config.OAuthPrefix)
}

func Test_getEnvConfigUseDefaults(t *testing.T) {
	config, err := getEnvConfig()
	require.Nil(t, err)
	require.Equal(t, true, config.MaxAuthEnabled)
	require.Equal(t, float64(1), config.MaxAuthRequestsPerSecond)
	require.Equal(t, 2, config.MaxAuthRequestBurst)
	require.Equal(t, false, config.HideDeprecated)
	require.Equal(t, "keptn:", config.OAuthPrefix)
	require.Equal(t, false, config.OAuthEnabled)
}

func TestPatchHtmlForOAuth(t *testing.T) {
	type args struct {
		readFile     func(string) ([]byte, error)
		oauthEnabled bool
		oauthPrefix  string
		deprecated   bool
	}
	html := `const oauth_prefix = "";
	const oauth_enabled = false;
	const hide_deprecated = false;`
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "reads from swagger-ui/index.html",
			args: args{
				readFile: func(s string) ([]byte, error) {
					var err error = nil
					if s == "swagger-ui/index.html" {
						err = errors.New("Must patch swagger-ui/index.html")
					}
					return []byte("ignore me due to error"), err
				},
				oauthEnabled: true,
				oauthPrefix:  "keptn:",
				deprecated:   true,
			},
			want: "",
		},
		{
			name: "do not patch",
			args: args{
				readFile:     func(s string) ([]byte, error) { return []byte(html), nil },
				oauthEnabled: false,
				oauthPrefix:  "keptn:",
				deprecated:   false,
			},
			want: html,
		},
		{
			name: "patch OAuth",
			args: args{
				readFile:     func(s string) ([]byte, error) { return []byte(html), nil },
				oauthEnabled: true,
				oauthPrefix:  "prefix:",
				deprecated:   false,
			},
			want: `const oauth_prefix = "prefix:";
			const oauth_enabled = true;
			const hide_deprecated = false;`,
		},
		{
			name: "patch Deprecated",
			args: args{
				readFile:     func(s string) ([]byte, error) { return []byte(html), nil },
				oauthEnabled: false,
				oauthPrefix:  "prefix:",
				deprecated:   true,
			},
			want: `const oauth_prefix = "";
			const oauth_enabled = false;
			const hide_deprecated = true;`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := getEnvConfig()
			require.Nil(t, err)
			config.HideDeprecated = tt.args.deprecated
			config.OAuthEnabled = tt.args.oauthEnabled
			config.OAuthPrefix = tt.args.oauthPrefix

			got := patchHtmlForOAuth(config, tt.args.readFile)
			require.Equal(t, strip(tt.want), strip(got))
		})
	}
}

func strip(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(in, "\t", ""), " ", "")
}

func TestHTTPHandler_normalIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/swagger-ui/index.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handlerAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handle(rr, req, "", handlerAPI)
	res, _ := ioutil.ReadAll(rr.Body)

	// we get redirected to the file
	require.Equal(t, rr.Code, 301)
	require.Equal(t, string(res), "")
}

func TestHTTPHandler_patchedIndex(t *testing.T) {
	content := "patched content"
	req, err := http.NewRequest("GET", "/swagger-ui/index.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handlerAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handle(rr, req, content, handlerAPI)
	res, _ := ioutil.ReadAll(rr.Body)

	require.Equal(t, rr.Code, 200)
	require.Equal(t, string(res), content)
}

func TestHTTPHandler_health(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	hasBeenCalled := false
	handlerAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasBeenCalled = true
	})
	handle(rr, req, "ignored", handlerAPI)
	res, _ := ioutil.ReadAll(rr.Body)

	require.Equal(t, rr.Code, 200)
	require.Equal(t, string(res), "")
	require.Equal(t, hasBeenCalled, false)
}

func TestHTTPHandler_api(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	hasBeenCalled := false
	handlerAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasBeenCalled = true
	})
	handle(rr, req, "ignored", handlerAPI)
	res, _ := ioutil.ReadAll(rr.Body)

	require.Equal(t, rr.Code, 200)
	require.Equal(t, string(res), "")
	require.Equal(t, hasBeenCalled, true)
}
