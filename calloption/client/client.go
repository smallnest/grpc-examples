package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/smallnest/grpc-examples/calloption/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	address = flag.String("addr", "localhost:8972", "address")
	name    = flag.String("n", "world", "name")
)

func main() {
	flag.Parse()

	// 连接服务器
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	//unary
	ctx := context.Background()

	p := &peer.Peer{}
	callOption := grpc.Peer(p)

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name}, callOption)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	log.Printf("peer: %+v", p)

}
