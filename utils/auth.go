package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

//Gin 验证
func GinAuth(c *gin.Context) {
	method := c.Request.Method
	//跨域
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	//token

	// 处理请求
	c.Next()
}

// grpc 验证
func GrpcAuth(ctx context.Context, info *grpc.UnaryServerInfo) error {
	fmt.Println(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return errors.New("Token认证信息无效")
	}

	return nil
}
