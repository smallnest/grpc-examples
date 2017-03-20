package main

import (
	"log"
	"net"

	"flag"

	pb "github.com/smallnest/grpc-examples/interceptor/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.String("p", ":8972", "port")
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(StreamServerInterceptor),
		grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterGreeterServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling. Info: %+v", info)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)
	return resp, err
}

func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("before handling. Info: %+v", info)
	err := handler(srv, ss)
	log.Printf("after handling. err: %v", err)
	return err
}
