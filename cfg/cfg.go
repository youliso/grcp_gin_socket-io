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
		Name:              "md1",
		Uri:               "root:pwd@tcp(127.0.0.1:3306)/db?parseTime=true&loc=Local&charset=utf8mb4",
		MysqlMaxOpenConns: 100,
		MysqlMaxIdleConns: 10,
	},
	{
		Ty:                  "redis",
		Name:                "md2",
		Uri:                 "redis://127.0.0.1:6379",
		Pwd:                 "",
		RedisMaxIdle:        3,                 //最大空闲连接数
		RedisIdleTimeoutSec: time.Second * 240, //最大空闲连接时间
	},
	//{
	//	Ty:             "mongo",
	//	Name:           "md3",
	//	Database:       "admin",
	//	Uri:            "127.0.0.1:27017",
	//	Pwd:            "",
	//	MonMaxPoolSize: 4096,
	//	MonTimeoutSec:  60 * time.Second,
	//},
}

type Dbt struct {
	Ty                  string
	Name                string
	Uri                 string
	Pwd                 string
	Database            string
	MysqlMaxOpenConns   int
	MysqlMaxIdleConns   int
	RedisMaxIdle        int
	RedisIdleTimeoutSec time.Duration
	MonMaxPoolSize      int
	MonTimeoutSec       time.Duration
}
