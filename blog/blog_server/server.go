package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var collection *mongo.Collection

type server struct {
}

type blogItem struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string	`bson:"author_id"`
	Content string	`bson:"content"`
	Title string `bson:"title"`
}

func (*server)CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error){
	fmt.Println("Create Blog Request")
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

	res, err := collection.InsertOne(nil, data)
	if err != nil{
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err),
			)
	}

	collid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can't convert collectionId"),
			)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: collid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func main(){
	// Set Log Flat
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to Database...")
	// Create Database Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Can't create Database client: %v", err)
	}

	// Connect to Database Server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Can't connect to Database: %v", err)
	}

	fmt.Println("Blog Server started...")
	collection = client.Database("mydbs").Collection("blog")

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
	_  = lis.Close()
	fmt.Println("Closing mongoDB Connection...")
	client.Disconnect(ctx)
	fmt.Println("End of Program...")
}