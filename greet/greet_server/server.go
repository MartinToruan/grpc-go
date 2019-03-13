package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error{
	for {
		req, err := stream.Recv()

		// Check EOF
		if err == io.EOF {
			fmt.Printf("Finish Got Request from client\n")
			break
		}

		// Check Error
		if err != nil{
			log.Fatalf("Error while reading request from client: %v", err)
			return err
		}

		// Return Response to Client
		if err := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + req.GetGreeting().GetFirstName(),
		}); err != nil{
			log.Fatalf("Error when send response to the client: %v", err)
			return err
		}
	}
	return nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline was invoke: %v\n", req)

	for i := 0; i < 3; i++{
		if ctx.Err() == context.Canceled{
			return nil, status.Errorf(codes.Canceled, "Client was cancelled the request, may be we need rollback transaction")
		}
		time.Sleep(1 * time.Second)
	}

	return &greetpb.GreetWithDeadlineResponse{
		Result: "Hello " + req.GetGreeting().GetFirstName(),
	}, nil
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
