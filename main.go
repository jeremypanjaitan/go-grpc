package main

import (
	"context"
	"flag"
	"fmt"
	pb "gprc-go/grpcgo"
	"log"
	"net"

	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedUserServer
}

func (s *server) CreateUser(ctx context.Context, in *pb.UserDataRequest) (*pb.UserCreatedReply, error) {
	log.Printf("Received: %v\n", in.GetName())
	return &pb.UserCreatedReply{Message: "user with name " + in.GetName() + " created", Data: in}, nil
}

type productServer struct {
	pb.UnimplementedProductServer
}

func (p *productServer) CreateProduct(ctx context.Context, in *pb.ProductDataRequest) (*pb.ProductCreatedReply, error) {
	log.Printf("Received: %v\n", in.GetName())
	return &pb.ProductCreatedReply{
		Message: "product with name " + in.GetName() + " created",
		Data:    in,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &server{})
	pb.RegisterProductServer(s, &productServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server : %v", err)
	}
}
