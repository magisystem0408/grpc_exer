package main

import (
	"context"
	"fmt"
	"go_grcp/calculator/calculatorpb"
	"google.golang.org/grpc"
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

	doUnary(c)
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
