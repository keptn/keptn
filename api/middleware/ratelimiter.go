package middleware

import (
	"github.com/benbjohnson/clock"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	RequestsPerSecond float64
	MaxBurstSize      int
	tokenValidator    TokenValidator
	theClock          clock.Clock
	visitors          map[string]*visitor
	mutex             *sync.Mutex
}

func NewRateLimiter(requestsPerSecond float64, maxBurstSize int, tokenValidator TokenValidator, theClock clock.Clock) *RateLimiter {
	rl := &RateLimiter{
		RequestsPerSecond: requestsPerSecond,
		MaxBurstSize:      maxBurstSize,
		theClock:          theClock,
		visitors:          map[string]*visitor{},
		mutex:             &sync.Mutex{},
		tokenValidator:    tokenValidator,
	}

	ticker := rl.theClock.Ticker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C
			rl.cleanIPBuckets()
		}
	}()

	return rl
}

func (r *RateLimiter) Handle(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.Apply(w, req, handler)
	})
}

func (r *RateLimiter) Apply(w http.ResponseWriter, req *http.Request, handler http.Handler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ipAddress := getRemoteIP(req)
	limiter := r.getIPBucket(ipAddress)
	if limiter.Allow() == false {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	_, err := r.tokenValidator.ValidateToken(req.Header.Get("x-token"))
	if err == nil {
		r.clearIPBucket(ipAddress)
	}
	handler.ServeHTTP(w, req)
}

func (r *RateLimiter) getIPBucket(ip string) *rate.Limiter {
	v, exists := r.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(r.RequestsPerSecond), r.MaxBurstSize)
		r.visitors[ip] = &visitor{limiter, r.theClock.Now().UTC()}
		return limiter
	}
	// Update the last seen time for the visitor.
	v.lastSeen = r.theClock.Now().UTC()
	return v.limiter
}

func (r *RateLimiter) clearIPBucket(ip string) {
	limiter := rate.NewLimiter(rate.Limit(r.RequestsPerSecond), r.MaxBurstSize)
	r.visitors[ip] = &visitor{limiter, r.theClock.Now().UTC()}
}

func (r *RateLimiter) cleanIPBuckets() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for ip, v := range r.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(r.visitors, ip)
		}
	}
}

func getRemoteIP(r *http.Request) string {
	// first, check 'x-real-ip'
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	// then, check 'x-forwarded-for' header
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		split := strings.Split(forwardedFor, ",")
		if len(split) > 0 {
			return strings.TrimSpace(split[0])
		}
	}
	// finally, try to use RemoteAddr
	if r.RemoteAddr != "" {
		s := extractIPFromRemoteAddress(r.RemoteAddr)
		if s != "" {
			return s
		}
	}

	return ""
}

func extractIPFromRemoteAddress(addr string) string {
	ip, _, err := net.SplitHostPort(addr)
	if err == nil && ip != "" {
		return ip
	}
	return addr
}
