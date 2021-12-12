package main

import (
	"context"
	"fmt"
	"go_grcp/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
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
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Recieved PrimeNumber RPC:%v \n", req)
	number := req.GetNumber()
	//ここにアルゴリズムを書く
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
			fmt.Println(divisor)
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Recieved ComputeAverage RPC")
	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avarage := float64(sum) / float64(count)

			//送信内容が切れたら実行される
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: avarage,
			})
		}
		if err != nil {
			log.Printf("error while reading client steram %v", err)
		}
		sum += req.GetNumber()
		count++
	}
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
