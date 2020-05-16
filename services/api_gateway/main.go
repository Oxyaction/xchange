package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/Oxyaction/xchange/rpc"
	"github.com/Oxyaction/xchange/services/api_gateway/routers"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkingrpcmiddleware "github.com/openzipkin/zipkin-go/middleware/grpc"
	zipkinhttpmiddleware "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	"google.golang.org/grpc"
)

const (
	address = "localhost:3000"
)

func main() {
	// Zipkin STARTED
	endpoint, err := zipkin.NewEndpoint("api_gateway", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	reporter := zipkinhttp.NewReporter("http://127.0.0.1:9411/api/v2/spans")

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	serverMiddleware := zipkinhttpmiddleware.NewServerMiddleware(
		tracer, zipkinhttpmiddleware.TagResponseSize(true),
	)

	client := zipkingrpcmiddleware.NewClientHandler(tracer)
	// Zipkin ENDED

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second), grpc.WithStatsHandler(client))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	accountClient := pb.NewAccountClient(conn)

	router := routers.InitRoutes(accountClient)
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMiddleware(router),
	}
	fmt.Println("Server started on port 8080")
	server.ListenAndServe()
}
