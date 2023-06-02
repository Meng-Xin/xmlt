package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
	"xmlt/global"
)

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	logFilePath := global.Config.LogSet.LogFilePath
	logFileName := global.Config.LogSet.LogFileName

	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	var src io.Writer
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		src, err = os.Create(fileName)
	} else {
		src, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由 截断，防止一些恶意攻击特别长的URI打爆日志
		maxURI := global.Config.LogSet.SaveMaxURI
		reqUri := c.Request.URL.String()
		if len(reqUri) > maxURI {
			reqUri = reqUri[:maxURI]
		}
		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s ",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
