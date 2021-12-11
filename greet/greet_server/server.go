package main

import (
	"context"
	"fmt"
	"go_grcp/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet funciton was invoked with %v", req)

	firstName := req.GetGreeting().GetFirstName()
	result := "hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("failded to listen %v", err)
	}

	s := grpc.NewServer()
	//どのサーバーにするかの指定
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failded to listen %v", err)
	}
}
