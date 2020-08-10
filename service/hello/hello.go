package hello

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	pb "grpc/protos/hello"
	"grpc/utils"
	"grpc/utils/db/mysql"
	"reflect"
	"time"
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
	rows, err := mysql.Query("md1", "select * from table", reflect.TypeOf(new(test)).Elem())
	if err != nil {
		println(err.Error())
	}
	c.JSON(200, gin.H{
		"message": c.Query("name"),
		"data":    rows,
	})
}

type test struct {
	Id           int       `column:"id"`
	JsTemplate   string    `column:"js_template"`
	Title        string    `column:"title"`
	Img          string    `column:"img"`
	Introduction string    `column:"introduction"`
	Html         string    `column:"html"`
	Vk           string    `column:"value_key"`
	Ct           time.Time `column:"creation_time"`
}
