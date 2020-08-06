package routers

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"grpc/dao/hello"
	pb "grpc/protos/hello"
)

func GrpcRouters(s *grpc.Server) {
	pb.RegisterHelloServer(s, &hello.HelloService{})
}

func GinRouters(r *gin.Engine) {
	r.GET("/http", hello.HttpHello)
}
