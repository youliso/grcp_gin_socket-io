package cfg

import (
	"time"
)

const (
	//监听地址
	Uri = "127.0.0.1:3000"
	//TLS 认证
	OpenTLS = false
)

var Db = []Dbt{
	{
		Ty:                "mysql",
		Uri:               "bf:a3yPhPKCY47c8tch@tcp(175.24.76.246:3306)/bf?parseTime=true&loc=Local&charset=utf8mb4",
		MysqlMaxOpenConns: 100,
		MysqlMaxIdleConns: 10,
	},
	{
		Ty:                  "redis",
		Uri:                 "redis://127.0.0.1:6379",
		Pwd:                 "",
		RedisMaxIdle:        3,                 //最大空闲连接数
		RedisIdleTimeoutSec: time.Second * 240, //最大空闲连接时间
	},
}

type Dbt struct {
	Ty                  string
	Uri                 string
	Pwd                 string
	MysqlMaxOpenConns   int
	MysqlMaxIdleConns   int
	RedisMaxIdle        int
	RedisIdleTimeoutSec time.Duration
}
