package controller

import (
	"github.com/gin-gonic/gin"
	"sync"
)

// GracefulShutdownMiddleware synchronize active handlers to enable graceful shutdown
func GracefulShutdownMiddleware(wg *sync.WaitGroup) gin.HandlerFunc {

	return func(c *gin.Context) {
		wg.Add(1)
		c.Next()
		wg.Done()
	}
}
