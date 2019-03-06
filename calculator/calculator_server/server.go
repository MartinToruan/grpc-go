package main

import (
	"fmt"
	"github.com/MartinToruan/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"net"
	"log"
	"context"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error){
	fmt.Printf("Got Request with Message: %v\n", req)
	return &calculatorpb.CalculatorResponse{
		Result: req.Calculator.X + req.Calculator.Y,
	}, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	value := req.GetValue()

	// Prime Iterator
	var k int32 = 2

	// Start Generate Prime Number Decomposition
	for {
		// Stop when value is 0
		if value == 0{
			break
		}

		if value % k == 0{
			value /= k

			// Push k to Client
			return stream.Send(&calculatorpb.PrimeNumberResponse{
				Result: k,
			})
		} else{
			k++
		}
	}
	return nil
}

func main(){
	fmt.Println("Server Started!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Can't Listen on Port 50051: %v\n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil{
		log.Fatalf("Can't Start the Server: %v\n", err)
	}
}
