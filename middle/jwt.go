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
			return
		}
		// 拿到并切分token
		tokens := strings.Split(author, ",")
		aToken, rToken, err := utils.RefreshToken(tokens[0], tokens[1])
		if err != nil {
			code = e.TokenExpired
			c.JSON(http.StatusOK, public.Response{Status: code, Msg: e.GetMsg(code)})
			c.Abort()
		}
		// 如果用户刷新Token，那么使用新Token 解析，如果没有到期使用老Token。
		var parseToken string
		if aToken != "" {
			// 通过响应头信息 无感通知客户端更新双Token
			parseToken = aToken
			c.Header("New-Token", aToken+","+rToken)
		} else {
			parseToken = tokens[0]
		}
		// 解析获取用户载荷信息
		payLoad, err := utils.VerifyToken(parseToken)
		if err != nil {
			return
		}
		c.Set("uid", payLoad.UserID)
		c.Set("username", payLoad.Username)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}
