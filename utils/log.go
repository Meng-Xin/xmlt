package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type ILog interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type LLogger struct {
	logger *logrus.Logger
}

type LogEmailHook struct {
}

// Levels 需要监控的日志等级，只有命中列表中的日志等级才会触发Hook
func (l *LogEmailHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

// Fire 触发钩子函数，本实例为触发后发送邮件报警。
func (l *LogEmailHook) Fire(entry *logrus.Entry) error {
	// 修改日志内容
	entry.Data["app"] = "email"
	// 发送邮件

	mailTo := []string{"wenhaoli2022@163.com"} // 目标邮箱
	subject := "日志系统：记录到致命错误"                  // 邮件主题
	msg, _ := entry.String()                   // 邮件正文
	// 发送邮件
	err := SendMail(mailTo, subject, msg)
	if err != nil {
		return err
	}
	return nil
}

func NewLLogger() ILog {
	logFile := "./config/log/systemLog.txt"
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile " + logFile)
		panic(err)
	}
	log := &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),         // 文件 + 控制台输出
		Level: logrus.InfoLevel,                     // Debug日志等级
		Hooks: make(map[logrus.Level][]logrus.Hook), // 初始化Hook Map,否则导致Hook添加过程中的空指针引用。
		Formatter: &logrus.TextFormatter{ // 文本格式输出
			FullTimestamp:   true,                  // 展示日期
			TimestampFormat: "2006-01-02 15:04:05", //日期格式
			ForceColors:     false,                 // 颜色日志
		},
	}
	log.AddHook(&LogEmailHook{})
	log.Infof("日志开启成功")
	return &LLogger{logger: log}
}
func (l *LLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
