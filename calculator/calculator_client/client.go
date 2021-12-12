package main

import (
	"context"
	"fmt"
	"go_grcp/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {

	fmt.Println("Caluculator Client")

	//withInscureはsslで通信するという意味
	//grcpをセキュアに行うオプション
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	//プログラムの最後にdeferでclientを閉める
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Println("Created client: %f", c)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

//Unaryでの実行結果
func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	//実行結果
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeNumberDecomposition Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Server Streaming RPC: %v", req)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happned: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a ComputeAverage client streaming RPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while opening stream:%v", err)
	}

	numbers := []int32{3, 5, 6, 1, 21, 4, 134}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v \n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response:%v", err)
	}

	fmt.Printf("the avarage is: %v \n ", res.GetAverage())

}
