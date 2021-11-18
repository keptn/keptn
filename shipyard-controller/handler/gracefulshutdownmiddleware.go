package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

// GracefulShutdownMiddleware syncronise active handlers to enable graceful shutdown
func GracefulShutdownMiddleware(wg *sync.WaitGroup) gin.HandlerFunc {

	return func(c *gin.Context) {
		wg.Add(1)
		ctx := context.WithValue(c.Request.Context(), gracefulShutdownKey, wg)
		c.Request.WithContext(ctx)
		c.Next()
		wg.Done()
	}
}
