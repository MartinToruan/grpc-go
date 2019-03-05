package main

import (
	"fmt"
	"github.com/MartinToruan/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main(){
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created Client: %f", c)




}
