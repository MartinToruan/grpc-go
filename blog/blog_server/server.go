package main

import (
	"context"
	"fmt"
	"github.com/MartinToruan/grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
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

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error){
	fmt.Println("Read Blog Request")

	blogId := req.GetBlogId()
	objId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil{
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
			)
	}

	// Create empty struct
	data := &blogItem{}

	filter := bson.M{"_id": objId}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil{
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error){
	fmt.Println("Update Blog Request")
	blog := req.GetBlog()
	objId, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil{
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprint("Cannot parse ID"))
	}

	// Create empty struct
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}
	filter := bson.M{"_id": objId}

	// Update data in Database
	_, err = collection.ReplaceOne(context.Background(), filter, data)
	if err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Can't update data in database: %v", err))
	}

	data.ID = objId
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(&data),
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error){
	fmt.Println("Delete Blog Request")
	blogId := req.GetBlogId()
	objId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil{
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID: %v", err))
	}

	// Create filter
	filter := bson.M{"_id": objId}

	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil{
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot found data in database: %v", err))
	}

	return &blogpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error{
	fmt.Println("List Blog Request")

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil{
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected Error: %v", err))
	}
	defer func(){
		_ = cur.Close(context.Background())
	}()

	for cur.Next(context.Background()){
		// Create Empty Struct
		data := &blogItem{}

		// Decode
		if err := cur.Decode(data); err != nil{
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Cannot parse data: %v", err))
		}

		if err := stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)}); err != nil{
			log.Fatalf("Unexpected Error: %v", err)
		}

	}

	if cur.Err() != nil{
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected Error: %v", cur.Err()))
	}

	return nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog{
	return &blogpb.Blog{
		Id: data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title: data.Title,
		Content: data.Content,
	}
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

	// Create Mirror
	reflection.Register(s)

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
	_ = client.Disconnect(ctx)
	fmt.Println("End of Program...")
}