package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Oxyaction/xchange/rpc"

	"google.golang.org/grpc"
)

const (
	address = "localhost:3000"
)

func getBalance(ctx context.Context, c pb.AccountClient) {
	reply, err := c.GetBalance(ctx, &pb.GetBalanceRequest{Id: "123e4567-e89b-12d3-a456-426655440000"})
	if err != nil {
		fmt.Printf("%+v", err)
	} else {
		fmt.Printf("Balance is %d \n", reply.GetBalance())
	}
}

func createAccount(ctx context.Context, c pb.AccountClient) {
	reply, err := c.Create(ctx, &pb.CreateRequest{})
	if err != nil {
		fmt.Printf("%+v", err)
	} else {
		fmt.Printf("New account id: '%s'\n", reply.GetId())
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := pb.NewAccountClient(conn)
	createAccount(ctx, c)
}
