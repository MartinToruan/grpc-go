package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
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

	//createBlog(c)

	//readBlog(c)

	//updateBlog(c)

	//deleteBlog(c)

	getBlogList(c)
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

func readBlog(c blogpb.BlogServiceClient){
	fmt.Println("Read Blog")
	req := &blogpb.ReadBlogRequest{
		BlogId: "5c908ead68bc2fc569ed31af",
	}

	// Call ReadBlog Service
	resp, err := c.ReadBlog(context.Background(), req)
	if err != nil{
		// Parse and Print Error
		respErr, ok := status.FromError(err)
		if ok{
			fmt.Println(respErr.Code())
			fmt.Println(respErr.Message())

			if respErr.Code() == codes.InvalidArgument{
				fmt.Println("Request is invalid.")
			}
			if respErr.Code() == codes.NotFound {
				fmt.Println("Data not Found.")
			}
		} else{
			log.Fatalf("Unexpected Error: %v", err)
		}
		return
	}
	fmt.Println("Response from server: %v", resp.GetBlog())
}

func updateBlog(c blogpb.BlogServiceClient){
	fmt.Println("Update Blog")
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id: "5c908ead68bc2fc569ed30af",
			AuthorId: "Martin Kristopel",
			Title: "Updated Field",
			Content: "This is the updated field",
		},
	}

	// Call UpdateBlog Service
	resp, err := c.UpdateBlog(context.Background(), req)
	if err != nil{
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Code())
			fmt.Println(respErr.Message())

			if respErr.Code() == codes.InvalidArgument{
				fmt.Println("Request is invalid.")
			}
			if respErr.Code() == codes.NotFound {
				fmt.Println("Data not Found.")
			}
		} else{
			log.Fatalf("Unexpected Error: %v", err)
		}
		return
	}
	fmt.Printf("Got response from the server: %v", resp.GetBlog())
}

func deleteBlog(c blogpb.BlogServiceClient){
	fmt.Println("Delete Blog")
	req := &blogpb.DeleteBlogRequest{
		BlogId: "5c908eb668bc2fc569ed30b0",
	}

	// Call Delete Blog Service
	resp, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		respStat, ok := status.FromError(err)
		if ok{
			fmt.Println(respStat.Code())
			fmt.Println(respStat.Message())
		} else{
			log.Fatalf("Unexpected Error: %v", err)
		}
		return
	}

	fmt.Printf("Got response from the server: %v", resp.GetBlogId())
}

func getBlogList(c blogpb.BlogServiceClient){
	fmt.Println("Get Blog List")

	resStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil{
		log.Fatalf("Error while invoke the server: %v", err)
	}
	for {
		resp, err := resStream.Recv()

		if err == io.EOF{
			break
		}

		if err != nil{
			log.Fatalf("Got an error response from server: %v", err)
		}
		fmt.Printf("Got response from the server: %v\n", resp.GetBlog())
	}
}