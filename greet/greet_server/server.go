package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	log.Printf("Server was invoked with message: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	return &greetpb.GreetResponse{
		Result: "Hai " + firstName,
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error{
	log.Printf("Server was invoked with message: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i:=0 ; i<10 ; i++{
		resp := &greetpb.GreetManyTimesResponse{
			Result: "Hi " + firstName + " counter: " + strconv.Itoa(i),
		}
		stream.Send(resp)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for{
		req, err := stream.Recv()
		if err == io.EOF{
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil{
			log.Fatalf("Error while receive request from client: %v", err)
			return err
		}
		fmt.Printf("Got Request with message: %v\n", req.GetGreeting())
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "!\n"
	}
}

func main(){
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil{
		log.Fatalf("Failed to server: %v", err)
	}
}
