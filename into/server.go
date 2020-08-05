package into

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc/cfg"
	"grpc/controller"
	gw "grpc/protos/hello"
	pb "grpc/protos/hello"
	"net"
	"net/http"
)

func Run() {
	var ins = make(chan int)
	num := 2
	go func() {
		Tpc()
		ins <- 1
	}()
	go func() {
		Http()
		ins <- 1
	}()
	for range ins {
		num--
		if num == 0 {
			close(ins)
		}
	}
}

func Tpc() {
	var creds credentials.TransportCredentials
	var s *grpc.Server
	var opts []grpc.ServerOption

	//拦截器
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("init")
		fmt.Println(info.FullMethod)

		err = controller.Auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if cfg.OpenTLS {
		// TLS认证
		creds, err = credentials.NewServerTLSFromFile("./cfg/keys/server.pem", "./cfg/keys/server.key")
		if err != nil {
			fmt.Errorf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s = grpc.NewServer(opts...)

	pb.RegisterHelloServer(s, &controller.HelloService{})
	fmt.Println("RPC Listen on " + cfg.Address)
	if err := s.Serve(lis); err != nil {
		fmt.Errorf("Failed to %v", err)
	}
}

func Http() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// grpc服务地址
	endpoint := cfg.Address
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// HTTP转grpc
	err := gw.RegisterHelloHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("HTTP Listen on " + cfg.AddressHttp)
	if err := http.ListenAndServe(cfg.AddressHttp, mux); err != nil {
		panic(err)
	}
}
