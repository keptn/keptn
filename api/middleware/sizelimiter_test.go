package middleware

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProvideHandlerThatLimitsSize(t *testing.T) {
	req, _ := http.NewRequest("POST", "/somepath", strings.NewReader("HELLO"))
	rr := httptest.NewRecorder()
	handler := EnforceMaxEventSize(1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		all, _ := ioutil.ReadAll(r.Body)
		require.Equal(t, "H", string(all))
	}))
	handler.ServeHTTP(rr, req)
}
