package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var rdb = make(map[string]*redis.Pool)

func InitRedis(DbName, uri, pwd string, redisMaxIdle int, redisIdleTimeoutSec time.Duration) {
	rdb[DbName] = &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: time.Second * redisIdleTimeoutSec,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(uri)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			if pwd != "" {
				if _, authErr := c.Do("AUTH", pwd); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func Set(DbName, k, v string) {
	c := rdb[DbName].Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func GetStringValue(DbName, k string) string {
	c := rdb[DbName].Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		fmt.Println("Get Error: ", err.Error())
		return ""
	}
	fmt.Println(username)
	return username
}

func SetKeyExpire(DbName, k string, ex int) {
	c := rdb[DbName].Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func CheckKey(DbName, k string) bool {
	c := rdb[DbName].Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return exist
	}
}

func DelKey(DbName, k string) error {
	c := rdb[DbName].Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SetJson(DbName, k string, data interface{}) error {
	c := rdb[DbName].Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func getJsonByte(DbName, k string) ([]byte, error) {
	c := rdb[DbName].Get()
	defer c.Close()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonGet, nil
}
