package db

import (
	"grpc/cfg"
	"grpc/utils/db/mysql"
	"grpc/utils/db/redis"
)

func Init() {
	for i := 0; i < len(cfg.Db); i++ {
		switch cfg.Db[i].Ty {
		case "mysql":
			mysql.InitMysql(cfg.Db[i].Name, cfg.Db[i].Uri, cfg.Db[i].MysqlMaxOpenConns, cfg.Db[i].MysqlMaxIdleConns)
		case "redis":
			redis.InitRedis(cfg.Db[i].Name, cfg.Db[i].Uri, cfg.Db[i].Pwd, cfg.Db[i].RedisMaxIdle, cfg.Db[i].RedisIdleTimeoutSec)
		}
	}

}
