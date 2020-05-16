package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Oxyaction/xchange/rpc"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"

	"github.com/openzipkin/zipkin-go"
	zipkingrpcmiddleware "github.com/openzipkin/zipkin-go/middleware/grpc"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const (
	port  = "3000"
	dbURL = "postgres://postgres:xchange@localhost/account?sslmode=disable&pool_max_conns=10"
)

func initDbPool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	return pool
}

func main() {
	endpoint, err := zipkin.NewEndpoint("api_gateway", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	reporter := zipkinhttp.NewReporter("http://127.0.0.1:9411/api/v2/spans")

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	statsHandler := zipkingrpcmiddleware.NewServerHandler(tracer)

	pool := initDbPool()
	defer pool.Close()
	fmt.Println("Connected to database")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(statsHandler))
	accountRepository := &accountRepository{pool}
	pb.RegisterAccountServer(s, &AccountServer{
		accountRepository: accountRepository,
	})

	fmt.Printf("Listening on port %s\n", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
