package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main(){
	log.Println("Hi! This is Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Can't connect to locahost:50051. Err: %v\n", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// do Call
	//doUnary(c)

	// Call PrimeNumberDecomposition
	doPrimeNumberDecomposition(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient){
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			X: 13,
			Y: 10,
		},
	}

	resp, err := c.Calculate(context.Background(), req)
	if err != nil{
		log.Fatalf("Error when trying to invoke API: %v\n", err)
	}
	fmt.Printf("Got Response from the server: %v\n", resp.Result)
}

func doPrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient){
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PrimeNumberRequest{
		Value: 120,
	})

	if err != nil{
		log.Fatalf("Failed when trying to call PrimeNumberDecomposition in server: %v\n", err)
	}

	fmt.Printf("Prime Number Decomposition of %d : ", 120)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil{
			log.Fatal("Error when get response from the server: %v", err)
		}
		fmt.Printf("%v ", res.Result)
	}
	fmt.Println()
}
