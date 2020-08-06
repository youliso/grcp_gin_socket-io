package db

import (
	"grpc/cfg"
)

func Init() {
	for i := 0; i < len(cfg.Db); i++ {
		switch cfg.Db[i].Ty {
		case "mysql":
			InitMysql(cfg.Db[i].Name, cfg.Db[i].Uri, cfg.Db[i].MysqlMaxOpenConns, cfg.Db[i].MysqlMaxIdleConns)
		case "redis":
			InitRedis(cfg.Db[i].Name, cfg.Db[i].Uri, cfg.Db[i].Pwd, cfg.Db[i].RedisMaxIdle, cfg.Db[i].RedisIdleTimeoutSec)
		}
	}
}

//mysql操作
//var as func(...db.Dba) *db.Query
//as = db.Table(md, "article")
//var ass []Article
//err = as().Select(&ass)
//if err != nil {
//	println(err.Error())
//}
//fmt.Println(ass)
