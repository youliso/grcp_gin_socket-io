package utils

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc/cfg"
)

// &pb.HelloRequest{Name: "gRPC"}
func Conn(addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if cfg.OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile("./cfg/keys/server.pem", "server name")
		if err != nil {
			fmt.Errorf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(addr, opts...)
}
