package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmlt/global"
)

func BucketRateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if global.TokenBucket.Allow() {
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
