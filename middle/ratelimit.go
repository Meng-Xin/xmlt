package middle

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func BucketRateLimiter(rate *rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rate.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  http.StatusTooManyRequests,
				"message": "too many request",
			})
			ctx.Abort()
			return
		}
	}
}
