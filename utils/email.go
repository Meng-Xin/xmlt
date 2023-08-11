package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"mime"
	"os"
)

// SendMail 发送邮件
func SendMail(mailTo []string, subject string, body string) error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": "***@qq.com",  //发送人邮箱（邮箱以自己的为准）
		"pass": "****",        //发送人邮箱的密码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
		"host": "smtp.qq.com", //邮箱服务器（此时用的是qq邮箱）
	}

	m := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "LLL")) // 添加别名
	m.SetHeader("To", mailTo...)                                  // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                               // 设置邮件主题
	m.SetBody("text/html", body)                                  // 设置邮件正文

	//一个文件（加入发送一个 txt 文件）：/tmp/foo.txt，我需要将这个文件以邮件附件的方式进行发送，同时指定附件名为：附件.txt
	//同时解决了文件名乱码问题
	name := "附件.txt"
	m.Attach("./gomail.txt",
		gomail.Rename(name), //重命名
		gomail.SetHeader(map[string][]string{
			"Content-Disposition": []string{
				fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", name)),
			},
		}),
	)

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(mailConn["host"], 465, mailConn["user"], mailConn["pass"]) // 设置邮件正文
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}

// ParseString 用来解析HTML
func ParseString(path string) string {
	byte, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("ReadFile error:", err)
		return ""
	}
	return string(byte)
}
