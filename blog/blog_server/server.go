package main

import (
	"fmt"
	"github.com/MartinToruan/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {

}

func main(){
	fmt.Println("Blog Server started...")

	// Set Log Flat
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create a Listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Error while starting up the server: %v", err)
	}

	// Create Grpc Server
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func(){
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil{
			log.Fatalf("Error while serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener...")
	lis.Close()
	fmt.Println("End of Program...")



}