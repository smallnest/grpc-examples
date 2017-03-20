package main

import (
	"context"
	"flag"
	"log"
	"strconv"

	pb "github.com/smallnest/grpc-examples/metatdata/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	md := metadata.Pairs("X-metadata-key1", "metadata-value1")
	ctx := metadata.NewContext(context.Background(), md)
	log.Printf("unary send MD: %+v", md)

	var header, trailer metadata.MD
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	log.Printf("Header: %s", header)
	log.Printf("Trailer: %s", trailer)

	//streaming
	md = metadata.Pairs("X-metadata-key2", "metadata-value2")
	ctx = metadata.NewContext(context.Background(), md)
	log.Printf("streaming send MD: %+v", md)

	stream, err := c.SayHello1(ctx)
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}

	for i := 0; i < 3; i++ {
		stream.Send(&pb.HelloRequest{Name: *name + strconv.Itoa(i)})
		if err != nil {
			log.Printf("failed to send: %v", err)
			break
		}
		reply, err := stream.Recv()
		if err != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("Greeting: %s", reply.Message)
	}
}
