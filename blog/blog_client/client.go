package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main(){
	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	createBlog(c)


}

func createBlog(c blogpb.BlogServiceClient){
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Kristopel Martin",
		Title: "My First Blog",
		Content: "Content of the first blog",
	}
	resp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil{
		log.Fatalf("Unexpected Error: %v", err)
	}

	fmt.Printf("Blog has been created: %v", resp)
}
