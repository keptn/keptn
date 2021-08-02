package middleware

import (
	"errors"
	"github.com/benbjohnson/clock"
	middleware_mock "github.com/keptn/keptn/api/middleware/fake"
	"github.com/keptn/keptn/api/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func Test_extractIPFromRemoteAddress(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "extract from ip4",
			args: args{
				addr: "127.0.0.1:8000",
			},
			want: "127.0.0.1",
		}, {
			name: "extract from ip6",
			args: args{
				addr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8000",
			},
			want: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
		{
			name: "extract from ip4 - without port",
			args: args{
				addr: "127.0.0.1",
			},
			want: "127.0.0.1",
		}, {
			name: "extract from ip6 - without port",
			args: args{
				addr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
			},
			want: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIPFromRemoteAddress(tt.args.addr); got != tt.want {
				t.Errorf("extractIPFromRemoteAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRemoteIPWithRemoteAddr(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)
	r.RemoteAddr = "127.0.0.1"

	ip := getRemoteIP(r)

	require.Equal(t, "127.0.0.1", ip)
}

func Test_getRemoteIPWithRemoteXForwardHeader(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)

	r.RemoteAddr = ""
	r.Header.Set("x-forwarded-for", "127.0.0.2,0.0.0.0")
	ip := getRemoteIP(r)

	require.Equal(t, "127.0.0.2", ip)
}

func Test_getRemoteIPWithRemoteXRealIP(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)

	r.RemoteAddr = ""
	r.Header.Set("x-real-ip", "127.0.0.2")
	ip := getRemoteIP(r)

	require.Equal(t, "127.0.0.2", ip)
}

func TestRateLimiter(t *testing.T) {
	mockClock := clock.NewMock()
	tokenValidator := &middleware_mock.TokenValidatorMock{ValidateTokenFunc: func(token string) (*models.Principal, error) {
		return nil, errors.New("oops")
	}}
	rl := NewRateLimiter(1.0, 1, tokenValidator, mockClock)

	mh := &MockHttpHandler{}
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.Nil(t, err)

	// send one request - this should pass the rate limiter
	rl.Apply(&httptest.ResponseRecorder{}, req, mh)
	require.Equal(t, 1, mh.calls)

	// send a couple of requests at once - these should not pass
	wg := &sync.WaitGroup{}
	wg.Add(6)

	for i := 0; i <= 5; i++ {
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodGet, "", nil)
			require.Nil(t, err)

			rl.Apply(&httptest.ResponseRecorder{}, req, mh)
		}()
	}
	wg.Wait()
	require.Equal(t, 1, mh.calls)

	// wait a bit - the next one should pass
	<-time.After(1 * time.Second)
	rl.Apply(&httptest.ResponseRecorder{}, req, mh)
	require.Equal(t, 2, mh.calls)

	// check if the visitor buckets are created
	require.Len(t, rl.visitors, 1)

	// proceed the internal clock of the limiter and check if the visitors are being cleaned up
	mockClock.Add(2 * time.Minute)
	require.Empty(t, rl.visitors)
}

type MockHttpHandler struct {
	calls int
	lock  sync.Mutex
}

func (mh *MockHttpHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	mh.lock.Lock()
	defer mh.lock.Unlock()
	mh.calls++
}
