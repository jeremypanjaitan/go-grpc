package main

import (
	"context"
	"flag"
	"fmt"
	pb "gprc-go/grpcgo"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	var values []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get("domain")
	}
	log.Println("Received metadata ", values[0])
	log.Printf("Received: %v\n", in.GetName())
	return &pb.ProductCreatedReply{
		Message: "product with name " + in.GetName() + " created",
		Data:    in,
	}, nil
}

func (p *productServer) GetBulkProduct(in *pb.GetBulkProductQuery, src pb.Product_GetBulkProductServer) error {
	var products []pb.ProductDataResponse

	// Create five product
	for i := 0; i < 5; i++ {
		products = append(products, pb.ProductDataResponse{
			Name:  fmt.Sprintf("Makanan %v", i+1),
			Price: in.GetPrice(),
		})
	}

	var wg sync.WaitGroup

	//Create 5 goroutine to send each of the product
	for i := 0; i < len(products); i++ {
		wg.Add(1)
		go func(count int, product *pb.ProductDataResponse) {
			defer wg.Done()
			time.Sleep(time.Duration(int64(1)) * time.Second)

			if err := src.Send(&pb.ProductBulkDataResponse{
				Message: fmt.Sprintf("data - %d", count),
				Data:    nil,
			}); err != nil {
				log.Println("disini", err)
			}
			log.Printf("Finishing request number : %d", count+1)
		}(i, &products[i])
	}

	//Wait until all goroutine finish executing
	wg.Wait()
	return nil
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var values []string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values = md.Get("domain")
		}
		if len(values) > 0 {
			log.Println("Received metadata from interceptor ", values[0])
		}
		return handler(ctx, req)
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor()))
	pb.RegisterUserServer(s, &server{})
	pb.RegisterProductServer(s, &productServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server : %v", err)
	}
}
