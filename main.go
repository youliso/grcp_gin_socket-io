package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc/cfg"
	"grpc/routers"
	"grpc/utils"
	"grpc/utils/db"
	"net"
	"net/http"
)

func main() {
	db.Init()
	run()
}

func run() {
	lis, err := net.Listen("tcp", cfg.Uri)
	if err != nil {
		println(err.Error())
		return
	}
	m := cmux.New(lis)
	httpl := m.Match(cmux.HTTP1Fast())
	grpcl := m.Match(cmux.Any())
	go gRpc(grpcl)
	go ginOrSocketIo(httpl)

	fmt.Println("Listen on " + cfg.Uri)
	if err := m.Serve(); err != nil {
		panic(err)
	}
}

func gRpc(lin net.Listener) {
	var creds credentials.TransportCredentials
	var s *grpc.Server
	var opts []grpc.ServerOption
	var err error
	//拦截器
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = utils.GrpcAuth(ctx, info)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	if cfg.OpenTLS {
		// TLS认证
		creds, err = credentials.NewServerTLSFromFile("./cfg/keys/server.pem", "./cfg/keys/server.key")
		if err != nil {
			panic(err.Error())
		}
		opts = append(opts, grpc.Creds(creds))
	}
	s = grpc.NewServer(opts...)
	routers.GrpcRouters(s)
	if err := s.Serve(lin); err != nil {
		panic(err.Error())
	}
}

func ginOrSocketIo(lin net.Listener) {
	r := gin.New()
	r.Use(utils.GinAuth)
	r.GET("/socket.io/*any", gin.WrapH(utils.SocketIo()))
	r.POST("/socket.io/*any", gin.WrapH(utils.SocketIo()))
	routers.GinRouters(r)
	s := &http.Server{
		Handler: r,
	}
	if err := s.Serve(lin); err != cmux.ErrListenerClosed {
		panic(err)
	}
}
