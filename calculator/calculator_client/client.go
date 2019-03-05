package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
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
	doUnary(c)
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
