package main

import (
	"context"
	"fmt"
	pb "go_grcp/pycli_goser/server/grpc-server"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct{}

func (*server) PushMsg(ctx context.Context, p *pb.MsgStruct) (*pb.MsgStruct, error) {
	log.Printf("Received:%v", p.Message)
	res := &pb.MsgStruct{Message: "hello" + p.Message}
	return res, nil
}

func main() {
	fmt.Println("サーバーが立ち上がりました。")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("faild to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMamushiServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("faild to serve: %v", err)
	}
}
