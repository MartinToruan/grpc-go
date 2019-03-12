package main

import (
	"fmt"
	"github.com/MartinToruan/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
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
	fmt.Printf("Got a request with message: %v\n", req)

	// Prime Iterator
	var k int32 = 2

	// Start Generate Prime Number Decomposition
	for {
		// Stop when value is 0
		if value <= 1{
			break
		}

		if value % k == 0{
			value /= k

			// Push k to Client
			if err := stream.Send(&calculatorpb.PrimeNumberResponse{
				Result: k,
			}); err != nil {
				return err
			}
		} else{
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error{
	var result float64
	counter := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result/float64(counter),
			})
		}
		fmt.Printf("Got request from Client: %v\n", req)
		if err != nil{
			log.Fatalf("Error when get request from client: %v", err)
			return err
		}
		val := req.GetValue()
		result += float64(val)
		counter++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	currentMax := int32(-999999)
	for{
		req, err := stream.Recv()

		// Check Client Finish Sending Request
		if err == io.EOF {
			fmt.Println("Finish Processing.")
			return nil
		}

		if err != nil {
			log.Fatalf("Error when get request from client: %v", err)
			return err
		}

		val := req.GetValue()
		fmt.Printf("Got Request from client: %v\n", val)
		if val > currentMax {
			currentMax = val
			if err:= stream.Send(&calculatorpb.FindMaximumResponse{
				Result: val,
			}); err != nil{
				log.Fatalf("Error when sending response to client: %v", err)
				return err
			}
		}
	}
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
