package main

import (
	"errors"
	"log"
	"net"

	"flag"

	pb "github.com/smallnest/grpc-examples/auth/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.String("p", ":8972", "port")
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if md, ok := metadata.FromContext(ctx); ok {
		key := md[":authority"]
		log.Printf("unary receive key: %+v", key)
		if key[0] != "a-secret-key" {
			return nil, errors.New("failed to auth")
		}
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
