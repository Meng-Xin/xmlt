package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/public"
	"xmlt/utils"
)

func VerifyJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		author := c.Request.Header.Get("Authorization")
		if len(author) == 0 {
			code = e.NotFindToken
			c.JSON(http.StatusOK, public.Response{Status: code, Msg: e.GetMsg(code)})
			c.Abort()
		}
		// 拿到并切分token
		tokens := strings.Split(author, ",")
		aToken, rToken, err := utils.RefreshToken(tokens[0], tokens[1])
		if err != nil {
			code = e.TokenExpired
			c.JSON(http.StatusOK, public.Response{Status: code, Msg: e.GetMsg(code)})
			c.Abort()
		}
		// 解析获取用户载荷信息
		payLoad, err := utils.VerifyToken(tokens[0])
		if err != nil {
			return
		}
		c.Set("uid", payLoad.ID)
		c.Set("username", payLoad.Username)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
		// 下面封装可以通知客户端更新双Token
		c.JSON(http.StatusOK, public.Response{Status: code, Data: aToken + "," + rToken, Msg: e.GetMsg(code)})
	}
}
