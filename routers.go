package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "grpc/protos/hello"
	"grpc/service/hello"
)

func GrpcRouters(s *grpc.Server) {
	pb.RegisterHelloServer(s, &hello.HelloService{})
}

func GinRouters(r *gin.Engine) {
	r.GET("/http", hello.HttpHello)
}
