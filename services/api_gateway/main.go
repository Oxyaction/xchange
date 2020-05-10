package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/Oxyaction/xchange/rpc"
	"github.com/Oxyaction/xchange/services/api_gateway/routers"

	"google.golang.org/grpc"
)

const (
	address = "localhost:3000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	accountClient := pb.NewAccountClient(conn)

	router := routers.InitRoutes(accountClient)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Server started on port 8080")
	server.ListenAndServe()
}
