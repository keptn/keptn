package handler

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type IAuthHandler interface {
	VerifyToken(c *gin.Context)
}

// type visitor struct {
// 	limiter  *rate.Limiter
// 	lastSeen time.Time
// }

type AuthHandler struct {
	// RequestsPerSecond float64
	// MaxBurstSize      int
	// theClock          clock.Clock
	// visitors          map[string]*visitor
	mutex *sync.Mutex
}

func NewAuthHandler( /*requestsPerSecond float64, maxBurstSize int, theClock clock.Clock*/ ) *AuthHandler {
	return &AuthHandler{
		// RequestsPerSecond: requestsPerSecond,
		// MaxBurstSize:      maxBurstSize,
		// theClock:          theClock,
		// visitors:          map[string]*visitor{},
		mutex: &sync.Mutex{},
	}

	// ticker := rl.theClock.Ticker(1 * time.Minute)
	// go func() {
	// 	for {
	// 		<-ticker.C
	// 		rl.cleanIPBuckets()
	// 	}
	// }()

	// return rl
}

// func (r *AuthHandler) VerifyToken(c *gin.Context) {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()
// 	ipAddress := getRemoteIP(c.Request)
// 	limiter := r.getIPBucket(ipAddress)
// 	if limiter.Allow() == false {
// 		SetTooManyRequestsErrorResponse(fmt.Errorf("too many requests"), c)
// 		return
// 	}
// 	_, err := r.validateToken(c.Request.Header.Get("x-token"))
// 	if err == nil {
// 		r.clearIPBucket(ipAddress)
// 	}

// 	c.JSON(http.StatusOK, nil)
// }

// func (r *AuthHandler) getIPBucket(ip string) *rate.Limiter {
// 	v, exists := r.visitors[ip]
// 	if !exists {
// 		limiter := rate.NewLimiter(rate.Limit(r.RequestsPerSecond), r.MaxBurstSize)
// 		r.visitors[ip] = &visitor{limiter, r.theClock.Now().UTC()}
// 		return limiter
// 	}
// 	// Update the last seen time for the visitor.
// 	v.lastSeen = r.theClock.Now().UTC()
// 	return v.limiter
// }

// func (r *AuthHandler) clearIPBucket(ip string) {
// 	limiter := rate.NewLimiter(rate.Limit(r.RequestsPerSecond), r.MaxBurstSize)
// 	r.visitors[ip] = &visitor{limiter, r.theClock.Now().UTC()}
// }

// func (r *AuthHandler) cleanIPBuckets() {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()
// 	for ip, v := range r.visitors {
// 		if time.Since(v.lastSeen) > 3*time.Minute {
// 			delete(r.visitors, ip)
// 		}
// 	}
// }

// func getRemoteIP(r *http.Request) string {
// 	// first, check 'x-real-ip'
// 	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
// 		return realIP
// 	}
// 	// then, check 'x-forwarded-for' header
// 	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
// 		split := strings.Split(forwardedFor, ",")
// 		if len(split) > 0 {
// 			return strings.TrimSpace(split[0])
// 		}
// 	}
// 	// finally, try to use RemoteAddr
// 	if r.RemoteAddr != "" {
// 		s := extractIPFromRemoteAddress(r.RemoteAddr)
// 		if s != "" {
// 			return s
// 		}
// 	}

// 	return ""
// }

// func extractIPFromRemoteAddress(addr string) string {
// 	ip, _, err := net.SplitHostPort(addr)
// 	if err == nil && ip != "" {
// 		return ip
// 	}
// 	return addr
// }

func (r *AuthHandler) VerifyToken(c *gin.Context) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	token := c.Request.Header.Get("x-token")
	if token == os.Getenv("SECRET_TOKEN") {
		c.JSON(http.StatusOK, nil)
		return
	}
	log.Errorf("Access attempt with incorrect api key auth: %s", token)
	SetUnauthorizedErrorResponse(fmt.Errorf("incorrect api key auth"), c)
}
