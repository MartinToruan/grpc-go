package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
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
	doStreamServer(c)
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
