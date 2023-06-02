package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitLoadConfig() *AllConfig {
	config := viper.New()
	// 设置读取路径
	config.AddConfigPath("./config")
	// 设置读取文件名字
	config.SetConfigName("config")
	// 设置读取文件类型
	config.SetConfigType("yaml")
	// 读取文件载体
	var configData *AllConfig
	// 读取配置文件
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Use Viper ReadInConfig Fatal error config err:%s \n", err))
	}
	// 查找对应配置文件
	err = config.Unmarshal(&configData)
	if err != nil {
		panic(fmt.Errorf("read config file to struct err: %s\n", err))
	}
	// 打印配置文件信息
	fmt.Println(configData)
	return configData
}

// AllConfig 整合Config
type AllConfig struct {
	Server      Server   `json:"server" yaml:"server"`
	Proxy       Proxy    `json:"proxy" yaml:"proxy"`
	ChatConn    ChatConn `json:"chat_conn" yaml:"chatConn"`
	LogSet      LogSet   `json:"log_set" yaml:"logSet"`
	MysqlMake   Mysql    `json:"mysql_make" yaml:"mysqlMake"`
	MysqlOnline Mysql    `json:"mysql_online" yaml:"mysqlOnline"`
}

type Server struct {
	Addr      string `json:"addr" yaml:"addr"`
	Port      string `json:"port" yaml:"port"`
	OpenProxy bool   `json:"open_proxy" yaml:"OpenProxy"`
}

func (s *Server) DSN() string {
	return s.Addr + ":" + s.Port
}

type Proxy struct {
	Addr string `json:"addr" yaml:"addr"`
	Port string `json:"port" yaml:"port"`
}

type ChatConn struct {
	Timeout         int64 `json:"timeout" yaml:"timeout"`                   // Second
	IdleConnTimeout int64 `json:"idle_conn_timeout" yaml:"idleConnTimeout"` // Hour
}

type LogSet struct {
	LogFilePath string `json:"log_file_path" yaml:"logFilePath"`
	LogFileName string `json:"log_file_name" yaml:"logFileName"`
	SaveMaxURI  int    `json:"save_max_uri" yaml:"saveMaxURI"`
}
