package hello

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

func GrpcHello() {
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

func HttpHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": c.Query("name"),
	})
}
