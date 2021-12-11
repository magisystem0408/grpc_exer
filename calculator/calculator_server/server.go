package main

import (
	"context"
	"fmt"
	"go_grcp/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	fmt.Printf("Recieved Sum RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	//サーバー側のロジックをここで書く

	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("failded to listen %v", err)
	}

	s := grpc.NewServer()
	//どのサーバーにするかの指定
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failded to listen %v", err)
	}
}
