package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	rl2 := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		if rl2.TakeAvailable(1) == 1 {
			c.Next()
			return
		}
		c.String(http.StatusOK, "rate limit...")
		c.Abort()
	}
}
