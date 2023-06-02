package config

type Mysql struct {
	Username     string `json:"username" yaml:"username"`             // 数据库用户名
	Password     string `json:"password" yaml:"password"`             // 数据库密码
	Dbname       string `json:"db-name" yaml:"db-name"`               // 数据库名
	Path         string `json:"path" yaml:"path"`                     // 服务器地址:端口
	Port         string `json:"port" yaml:"port"`                     //:端口
	AutoMigrate  bool   `json:"auto_migrate" yaml:"autoMigrate"`      // 是否开启表迁移
	Config       string `json:"config" yaml:"config"`                 // 高级配置
	MaxIdleConns int    `json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	LogMode      string `json:"log-mode" yaml:"log-mode"`             // 是否开启Gorm全局日志
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
