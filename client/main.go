package main

import (
	"context"
	pb "gprc-go/grpcgo"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr = "localhost:50051"
)

func main() {

	//try to connect to grpc
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}
	defer conn.Close()

	//Create grpc client using the connection that has been established
	c := pb.NewUserClient(conn)

	//create timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//call the server function
	r, err := c.CreateUser(ctx, &pb.UserDataRequest{
		Name:    "Jeremy Panjaitan",
		Address: "Jakarta",
		Age:     22,
	})
	if err != nil {
		log.Fatalf("could not create user : %v", err)
	}

	log.Printf("result: %s", r.GetMessage())

	d := pb.NewProductClient(conn)
	rr, err := d.CreateProduct(ctx, &pb.ProductDataRequest{
		Name:  "Makanan",
		Price: 100,
	})
	if err != nil {
		log.Fatalf("could not create user : %v", err)
	}
	log.Printf("result: %s", rr.GetMessage())

}
