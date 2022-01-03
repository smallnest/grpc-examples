package main

import (
	"flag"
	"io"
	"log"
	"net"

	pb "github.com/smallnest/grpc-examples/calloption/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

var port = flag.String("p", ":8972", "port")

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("unary receive Peer: %+v", p)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("unary receive MD: %+v", md)
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHello1(gs pb.Greeter_SayHello1Server) error {
	if md, ok := metadata.FromIncomingContext(gs.Context()); ok {
		log.Printf("streaming receive MD: %+v", md)
	}

	for {
		in, err := gs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		gs.Send(&pb.HelloReply{Message: "Hello " + in.Name})
	}

	return nil
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
