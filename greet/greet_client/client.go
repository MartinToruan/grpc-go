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
	"time"
)

func main(){
	fmt.Println("Hello I'm a client\n")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created Client: %f", c)
	//doUnary(c)
	//doStreamServer(c)
	//doClientStream(c)
	//doBiDiStream(c)
	doUnaryWithDeadline(c)
}

func doBiDiStream(c greetpb.GreetServiceClient){
	stream, err := c.GreetEveryone(context.Background())
	if err != nil{
		log.Fatalf("Error while invoke the server: %v", err)
	}

	// Initialize variable
	chanDone := make(chan struct{})

	// Send Request to the server
	go func(){

		// Define lIst of Name
		var nameRequest = make(map[string]string)
		nameRequest["Martin"] = "Kristopel"
		nameRequest["Markus"] = "Erikson"
		nameRequest["Frans"] = "Kristian"
		for firstName, lastName := range nameRequest{
			err := stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: firstName,
					LastName: lastName,
				},
			})

			if err != nil{
				log.Fatalf("Error while send request to the server: %v", err)
			}
		}
		if err := stream.CloseSend(); err != nil{
			log.Fatalf("Error while Close stream to server: %v", err)
		}
	}()

	// Get Response from the server
	go func(){
		for{
			defer close(chanDone)
			resp, err := stream.Recv()

			// Check EOF
			if err == io.EOF{
				break
			}

			if err != nil{
				log.Fatalf("Error while reading reponse from the server: %v", err)
				return
			}
			fmt.Printf("Response from the server: %v\n", resp.Result)
		}
	}()

	// Wait until process Done
	<-chanDone
	fmt.Println("Finish Greet EveryOne")

}

func doClientStream(c greetpb.GreetServiceClient){
	stream, err := c.LongGreet(context.Background())
	if err != nil{
		log.Fatalf("Error while invoke the server: %v", err)
		return
	}

	// Define List of Name
	var nameRequest = make(map[string]string)
	nameRequest["Martin"] = "Kristopel"
	nameRequest["Markus"] = "Erikson"
	nameRequest["Frans"] = "Kristian"
	for firstName, lastName := range nameRequest{
		stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: firstName,
				LastName: lastName,
			},
		})
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil{
		log.Fatalf("Error while get response from the server: %v", err)
	}
	fmt.Printf("Got Response from server: \n%v\n", resp.Result)
}

func doUnary(c greetpb.GreetServiceClient){
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kristopel",
			LastName: "Martin",
		},
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil{
		log.Fatalf("Can't invoke the server: %v\n", err)
	}
	log.Printf("Got Response from the server: %v\n", resp.Result)
}

func doStreamServer(c greetpb.GreetServiceClient){
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kristopel",
			LastName: "Martin",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil{
		log.Fatalf("Can't invoke GreetManyTimes: %v", err)
	}
	for {
		resp, err := stream.Recv()

		if err == io.EOF{
			break
		}
		if err != nil{
			log.Fatalf("Error while invoke GreetManyTimes: %v", err)
		}
		fmt.Printf("Got Response from GreetManyTimes: %v\n", resp.Result)
	}
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration){
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Martin",
			LastName: "Toruan",
		},
	}

	// Create context with Deadline
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := c.GreetWithDeadline(ctx, req)

	if err != nil{
		greetErr, ok := status.FromError(err)

		if ok{
			if greetErr.Code() == codes.DeadlineExceeded {
				log.Fatalf("Timeout hits!")
			} else {
				log.Fatalf("Error while hit GreetWithDeadline: %v", greetErr)
			}
		} else{
			log.Fatalf("Unknown Error: %v", err)
		}

		return
	}
	log.Printf("Got reponse from the server: %v\n", resp.GetResult())
}
