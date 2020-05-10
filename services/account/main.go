package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Oxyaction/xchange/rpc"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
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
	pool := initDbPool()
	defer pool.Close()
	fmt.Println("Connected to database")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	accountRepository := &accountRepository{pool}
	pb.RegisterAccountServer(s, &AccountServer{
		accountRepository: accountRepository,
	})

	fmt.Printf("Listening on port %s\n", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
