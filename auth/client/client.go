package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/smallnest/grpc-examples/auth/pb"
	"google.golang.org/grpc"
)

var (
	address = flag.String("addr", "localhost:8972", "address")
	name    = flag.String("n", "world", "name")
)

func main() {
	flag.Parse()

	// 连接服务器
	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithAuthority("a-secret-key"))
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
