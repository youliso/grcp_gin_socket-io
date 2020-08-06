package controller

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "grpc/protos/hello"
	"grpc/utils"
)

// server is used to implement helloworld.GreeterServer.
type HelloService struct{}

// SayHello implements helloworld.GreeterServer
func (h *HelloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md)
	}

	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func Hello() {
	conn, err := utils.Conn("127.0.0.1:3000")
	if err != nil {
		println(err.Error())
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)
	res, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "123"})
	if err != nil {
		println(err.Error())
	}
	fmt.Println(res)
}

// auth 验证Token
func Auth(ctx context.Context, info *grpc.UnaryServerInfo) error {
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
