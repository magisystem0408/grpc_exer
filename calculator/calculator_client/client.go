package main

import (
	"context"
	"fmt"
	"go_grcp/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
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
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c, -2)
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream and calling findMaximum:%v", err)
	}

	waitc := make(chan struct{})
	go func() {
		numbers := []int32{4, 1, 6, 7, 1234, 3131313, 3131}
		for _, number := range numbers {
			fmt.Printf("Sending number...:%v", number)
			stream.Send(&calculatorpb.FindMaximumRequest{Number: number})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Recieved a new maximum of... :%v \n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient, n int32) {
	fmt.Println("Starting to do a square Unary Response RPC...")
	//	correct code
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("we probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error calling squareroller: %v", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v:%v \n", n, res.GetNumberRoot())
}
