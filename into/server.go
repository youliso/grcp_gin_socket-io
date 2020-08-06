package into

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc/cfg"
	"grpc/controller"
	pb "grpc/protos/hello"
	"grpc/utils"
	"net"
	"net/http"
)

func Run() {
	lis, err := net.Listen("tcp", cfg.Uri)
	if err != nil {
		println(err.Error())
		return
	}
	m := cmux.New(lis)
	httpl := m.Match(cmux.HTTP1Fast())
	grpcl := m.Match(cmux.Any())
	go Tpc(grpcl)
	go GinOrSocketIo(httpl)

	fmt.Println("Listen on " + cfg.Uri)
	if err := m.Serve(); err != nil {
		panic(err)
	}
}

func Tpc(lin net.Listener) {
	var creds credentials.TransportCredentials
	var s *grpc.Server
	var opts []grpc.ServerOption
	var err error
	//拦截器
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = controller.Auth(ctx, info)
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
	pb.RegisterHelloServer(s, &controller.HelloService{})
	if err := s.Serve(lin); err != nil {
		panic(err.Error())
	}
}

func GinOrSocketIo(lin net.Listener) {
	r := gin.New()
	r.Use(utils.Cors())
	io, err := socketio.NewServer(nil)
	if err != nil {
		panic(err.Error())
	}
	//io.set('transports', ['websocket', 'xhr-polling', 'jsonp-polling', 'htmlfile', 'flashsocket']);
	io.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	io.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	io.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	io.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	io.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go io.Serve()
	defer io.Close()

	r.GET("/socket.io/*any", gin.WrapH(io))
	r.POST("/socket.io/*any", gin.WrapH(io))

	r.GET("/http", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": c.Query("name"),
		})
	})
	s := &http.Server{
		Handler: r,
	}
	if err := s.Serve(lin); err != cmux.ErrListenerClosed {
		panic(err)
	}
}
