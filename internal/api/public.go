package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PublicController interface {
	Ping()
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Pong...")
}

func NotFount(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, "NotFound")
}

func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, "Unauthorized")
}
