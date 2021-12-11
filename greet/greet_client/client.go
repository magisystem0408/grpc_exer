package main

import (
	"context"
	"fmt"
	"go_grcp/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {

	fmt.Println("Hello, I'm a client")

	//withInscureはsslで通信するという意味
	//grcpをセキュアに行うオプション
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	//プログラムの最後にdeferでclientを閉める
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Println("Created client: %f", c)

	doUnary(c)
}

//Unaryでの実行結果
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "nekomamushi",
			LastName:  "aaaa"},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	//実行結果
	log.Printf("Response from Great: %v", res.Result)

}
