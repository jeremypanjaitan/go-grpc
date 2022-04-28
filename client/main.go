package main

import (
	"context"
	pb "gprc-go/grpcgo"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	addr = "localhost:50051"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

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
	createProductCtx := metadata.AppendToOutgoingContext(ctx, "domain", "example.com")
	rr, err := d.CreateProduct(createProductCtx, &pb.ProductDataRequest{
		Name:  "Makanan",
		Price: 100,
	})
	if err != nil {
		log.Fatalf("could not create user : %v", err)
	}
	log.Printf("result: %s", rr.GetMessage())

	stream, err := d.GetBulkProduct(context.Background(), &pb.GetBulkProductQuery{
		Price: 2000,
	})
	if err != nil {
		log.Println(err)
	}
	done := make(chan bool)

	//create go routine until that receive response sent by
	//the server until there is no response anymore
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("Resp received: %s", resp.GetMessage())
		}
	}()

	//wait until for the goroutine finish exeucting

	<-done
	log.Printf("finished")

}
