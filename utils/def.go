package utils

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
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

func SocketIo() *socketio.Server {
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

	return io
}
