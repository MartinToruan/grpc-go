package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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
	//doPrimeNumberDecomposition(c)

	// Call ComputeAverage service
	//doComputeAverage(c)

	// Call FindMaximum service
	doFindMaximum(c)
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

func doComputeAverage(c calculatorpb.CalculatorServiceClient){
	stream, err := c.ComputeAverage(context.Background())
	if err != nil{
		log.Fatalf("Error while call Compute Average Service: %v\n", err)
	}

	for i:=1;i<=4;i++{
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Value: int32(i),
		})
		if err != nil {
			log.Fatalf("Error while send request to Server: %v\n", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil{
		log.Fatalf("Got Response Error from the server: %v\n", err)
	}
	fmt.Printf("Got Response from the server: %v\n", resp.Result)
}

func doFindMaximum(c calculatorpb.CalculatorServiceClient){
	// create doneChannel
	doneChannel := make(chan struct{})

	// Create Client Stream
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while connect to the server: %v", err)
	}

	// Sending Stream Request
	go func(){
		listNumber := []int32{1, 5, 3, 6, 2, 20}
		for _, v := range listNumber {
			if err := stream.Send(&calculatorpb.FindMaximumRequest{
				Value: v,
			}); err != nil{
				log.Fatalf("Error while sending request to the server: %v", err)
			}
			time.Sleep(3 * time.Second)
		}
		if err := stream.CloseSend(); err != nil{
			log.Fatalf("Error while close close send stream to server: %v", err)
		}
	}()

	// Receive Stream Response
	go func(){
		for{
			resp, err := stream.Recv()

			// Check if Server finish Sending Response
			if err == io.EOF{
				break
			}

			if err != nil{
				log.Fatalf("Error while receive response from the server: %v", err)
				break
			}

			res := resp.GetResult()
			fmt.Printf("Got Response from server: %v\n", res)
		}
		close(doneChannel)
	}()

	<-doneChannel
	fmt.Println("Finish Process Response from server.")

}